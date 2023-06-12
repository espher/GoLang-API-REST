package sqs

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

// CreateAWSSession crea y retorna una sesión de AWS.
func CreateAWSSession() (*session.Session, error) {
	creds := credentials.NewStaticCredentials(
		os.Getenv("AWS_ACCESS_KEY_ID"),
		os.Getenv("AWS_SECRET_ACCESS_KEY"),
		"",
	)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")), // Reemplaza con la región deseada
		Credentials: creds,
	})
	if err != nil {
		return nil, err
	}

	return sess, nil
}
