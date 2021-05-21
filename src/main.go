//+build !test

package main

import (
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

// The main function will initiate the Lambda handler.
func main() {
	region := os.Getenv("AWS_REGION")
	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region)},
	))
	secretsService := secretsmanager.New(session)
	service := lambdaService{
		SecretsManager: secretsService,
		Session:        session,
	}

	lambda.Start(service.LambdaHandler)
}
