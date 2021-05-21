package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

type stepRotator struct {
	ClientRequestToken string
	SecretID           string
	SecretsManager     *secretsmanager.SecretsManager
}

func (s *stepRotator) CreateSecret() error {
	var (
		currentCredentials databaseCredentials
		err                error
		pendingCredentials databaseCredentials
	)

	// Try getting the pending secret, and if that fails, then get the current
	// secret.
	if pendingCredentials, err = s.fetchDatabaseCredentials(AWSPENDING); err == nil {
		log.Printf("CreateSecret: Successfully got pending secret\n")
	} else if currentCredentials, err = s.fetchDatabaseCredentials(AWSCURRENT); err != nil {
		err = fmt.Errorf("CreateSecret: Failed to get current secret: %w", err)
		log.Printf("%v\n", err)
		return err
	} else {
		// If pending version is unavailable and current version is available,
		// set random password.
		log.Printf("CreateSecret: Successfully got current secret\n")

		pendingCredentials = currentCredentials

		if pendingCredentials.Password, err = s.generateNewPassword(); err != nil {
			err = fmt.Errorf("CreateSecret: Could not generate new password: %w", err)
			log.Printf("%v\n", err)
			return err
		}

		pendingValueVersionID := s.saveDatabaseCredentials(pendingCredentials, AWSPENDING)

		log.Printf("CreateSecret: Successfully put pending secret with version ID: %s\n", pendingValueVersionID)
	}

	return nil
}

func (s *stepRotator) fetchDatabaseCredentials(stage versionStage) (databaseCredentials, error) {
	var credentials databaseCredentials

	versionStageStr := stage.String()

	input := &secretsmanager.GetSecretValueInput{
		SecretId:     &s.SecretID,
		VersionStage: &versionStageStr,
	}
	secretOutput, err := s.SecretsManager.GetSecretValue(input)

	if err != nil {
		err = fmt.Errorf("Could not retrieve secret value: %w", err)
		log.Printf("%v\n", err)
		return credentials, err
	}

	credentialsJSON := secretOutput.SecretString

	if credentialsJSON == nil {
		err = errors.New("Credentials are empty")
		log.Printf("%v\n", err)
		return credentials, err
	}

	if err = json.Unmarshal([]byte(*credentialsJSON), &credentials); err != nil {
		err = fmt.Errorf("Could not unmarshal credentials: %w", err)
		log.Printf("%v\n", err)
		return credentials, err
	}

	return credentials, nil
}

func (s *stepRotator) FinishSecret() error {
	var (
		currentVersionID string
		err              error
		pendingVersionID string
	)

	log.Printf("(FinishSecret) Start finishing secret for secret ID %s...\n", s.SecretID)

	if currentVersionID, err = s.getVersionIDFromVersionStage(AWSCURRENT); err != nil {
		return err
	}

	if pendingVersionID, err = s.getVersionIDFromVersionStage(AWSPENDING); err != nil {
		return err
	}

	if currentVersionID == s.ClientRequestToken {
		log.Printf(
			"(FinishSecret) Secret version with client request token already set as AWSCURRENT for version ID: %s\n",
			currentVersionID,
		)
	} else if err = s.moveVersionStage(AWSCURRENT, currentVersionID, pendingVersionID); err != nil {
		err = fmt.Errorf("Failed to update version stage: %w", err)
		log.Printf("%v\n", err)
		return err
	} else {
		log.Printf("(FinishSecret) Successfully updated AWSCURRENT from %#v to %#v.\n", currentVersionID, pendingVersionID)
	}

	return nil
}

func (s *stepRotator) generateNewPassword() (string, error) {
	var (
		err    error
		output *secretsmanager.GetRandomPasswordOutput
	)

	input := &secretsmanager.GetRandomPasswordInput{
		ExcludeLowercase:        aws.Bool(false),
		ExcludeNumbers:          aws.Bool(false),
		ExcludePunctuation:      aws.Bool(true),
		ExcludeUppercase:        aws.Bool(false),
		IncludeSpace:            aws.Bool(false),
		PasswordLength:          aws.Int64(32),
		RequireEachIncludedType: aws.Bool(true),
	}

	if output, err = s.SecretsManager.GetRandomPassword(input); err != nil {
		log.Fatalf("Failed to generate a random password; %s", err.Error())
		return "", err
	}

	return *output.RandomPassword, nil
}

// getVersionIDFromVersionStage returns a specific version ID for a stage.
func (s *stepRotator) getVersionIDFromVersionStage(vStage versionStage) (string, error) {
	versionIDsToStages, err := s.getVersionIDsToStages()

	if err != nil {
		return "", err
	}

	for id, stages := range versionIDsToStages {
		for _, stage := range stages {
			if *stage == string(vStage) {
				return id, nil
			}
		}
	}

	return "", nil
}

func (s *stepRotator) getVersionIDsToStages() (map[string][]*string, error) {
	var (
		err    error
		output *secretsmanager.DescribeSecretOutput
	)

	input := &secretsmanager.DescribeSecretInput{SecretId: &s.SecretID}

	if output, err = s.SecretsManager.DescribeSecret(input); err != nil {
		return map[string][]*string{},
			fmt.Errorf("Failed to describe secret with secret ID %s: %w", s.SecretID, err)
	}

	return output.VersionIdsToStages, nil
}

// moveVersionStage moves a version stage tag across versions.
func (s *stepRotator) moveVersionStage(stage versionStage, fromVerID string, toVerID string) error {
	verStgStr := stage.String()
	input := &secretsmanager.UpdateSecretVersionStageInput{
		MoveToVersionId:     &toVerID,
		RemoveFromVersionId: &fromVerID,
		SecretId:            &s.SecretID,
		VersionStage:        &verStgStr,
	}

	if _, err := s.SecretsManager.UpdateSecretVersionStage(input); err != nil {
		err = fmt.Errorf("Failed to move version stage %s from %#v to %#v: %w", stage.String(), fromVerID, toVerID, err)
		log.Printf("%v\n", err)
		return err
	}

	return nil
}

// saveDatabaseCredentials saves the database credentials in Secrets Manager.
func (s *stepRotator) saveDatabaseCredentials(credentials databaseCredentials, stage versionStage) string {
	var (
		output    *secretsmanager.PutSecretValueOutput
		versionID string
	)

	verStgStr := stage.String()
	bytes, err := json.Marshal(credentials)
	secret := string(bytes)
	input := &secretsmanager.PutSecretValueInput{
		ClientRequestToken: &s.ClientRequestToken,
		SecretId:           &s.SecretID,
		SecretString:       &secret,
		VersionStages:      []*string{&verStgStr},
	}

	if err != nil {
		log.Fatalf("Failed to marshal secret value; %s\n", err.Error())
	} else if output, err = s.SecretsManager.PutSecretValue(input); err != nil {
		log.Fatalf("Failed to put secret value; %s\n", err.Error())
	} else {
		versionID = *output.VersionId
	}

	return versionID
}

func (s *stepRotator) SetSecret() error {
	var (
		currentCredentials databaseCredentials
		err                error
		pendingCredentials databaseCredentials
	)

	log.Printf("(SetSecret) Start setting secret for secret ID %s...\n", s.SecretID)

	if currentCredentials, err = s.fetchDatabaseCredentials(AWSCURRENT); err != nil {
		err = fmt.Errorf("(SetSecret) Failed to get current secret: %w", err)
		log.Printf("%v\n", err)
		return err
	} else if pendingCredentials, err = s.fetchDatabaseCredentials(AWSPENDING); err != nil {
		err = fmt.Errorf("(SetSecret) Failed to get pending secret: %w", err)
		log.Printf("%v\n", err)
		return err
	}

	rotator, ok := engineRotators[currentCredentials.Engine]

	if !ok {
		err = fmt.Errorf("Unsupported database engine: %s", currentCredentials.Engine)
		log.Printf("%v\n", err)
		return err
	}

	if err = rotator.Rotate(currentCredentials, pendingCredentials.Password); err != nil {
		log.Printf("Could not rotate password: %v\n", err)
		return err
	}

	log.Printf("(SetSecret) Successfully updated credentials\n")

	return nil
}

func (s *stepRotator) TestSecret() error {
	var (
		err                error
		pendingCredentials databaseCredentials
	)

	log.Printf("(TestSecret) Testing secret ID %s...\n", s.SecretID)

	if pendingCredentials, err = s.fetchDatabaseCredentials(AWSPENDING); err != nil {
		err = fmt.Errorf("(TestSecret) Failed to get pending secret: %w", err)
		log.Printf("%v\n", err)
		return err
	}

	rotator, ok := engineRotators[pendingCredentials.Engine]

	if !ok {
		err = fmt.Errorf("Unsupported database engine: %s", pendingCredentials.Engine)
		log.Printf("%v\n", err)
		return err
	}

	if err = rotator.Test(pendingCredentials); err != nil {
		log.Printf("Could not test credentials: %v\n", err)
		return err
	}

	log.Printf("(TestSecret) Successfully tested credentials\n")

	return nil
}
