package main

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// rotationStep is a secret rotation step.
type rotationStep string

// versionStage is a Secret Manager Version ID.
type versionStage string

// lambdaEvent is the event our handler will receive.
type lambdaEvent struct {
	ClientRequestToken string       `json:"ClientRequestToken"`
	SecretID           string       `json:"SecretId"`
	Step               rotationStep `json:"Step"`
}

type lambdaService struct {
	SecretsManager *secretsmanager.SecretsManager
	Session        *session.Session
}

const (
	// ROTATION STEPS

	// CreateSecret specifies the secret creation step.
	createSecret rotationStep = "createSecret"
	// SetSecret specifies the secret setting step.
	setSecret rotationStep = "setSecret"
	// TestSecret specifies the secret testing step.
	testSecret rotationStep = "testSecret"
	// FinishSecret specifies the secret finishing step.
	finishSecret rotationStep = "finishSecret"

	// VERSION STAGES

	// AWSCURRENT is a version stage tag.
	AWSCURRENT versionStage = "AWSCURRENT"
	// AWSPENDING is a version stage tag.
	AWSPENDING versionStage = "AWSPENDING"
	// AWSPREVIOUS is a version stage tag.
	AWSPREVIOUS versionStage = "AWSPREVIOUS"
)

func (v versionStage) String() string {
	return string(v)
}

func (l *lambdaService) LambdaHandler(ctx context.Context, event lambdaEvent) error {
	var err error

	region := os.Getenv("AWS_REGION")
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)},
	))
	secretsService := secretsmanager.New(session)
	rotator := stepRotator{
		ClientRequestToken: event.ClientRequestToken,
		SecretID:           event.SecretID,
		SecretsManager:     secretsService,
	}

	switch event.Step {
	case createSecret:
		err = rotator.CreateSecret()
	case setSecret:
		err = rotator.SetSecret()
	case testSecret:
		err = rotator.TestSecret()
	case finishSecret:
		err = rotator.FinishSecret()
	}

	if err != nil {
		log.Printf("Final error: %v\n", err)
	}

	return err
}
