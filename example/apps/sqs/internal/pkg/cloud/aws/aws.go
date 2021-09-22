package aws

import (
"github.com/aws/aws-sdk-go/aws"
"github.com/aws/aws-sdk-go/aws/credentials"
"github.com/aws/aws-sdk-go/aws/session"
)

type Config struct {
	Address string
	Region  string
	Profile string
	ID      string
	Secret  string
}

func New(config Config) (*session.Session, error) {
	return session.NewSessionWithOptions(
		session.Options{
			Config: aws.Config{
				Credentials:      credentials.NewStaticCredentials(config.ID, config.Secret, ""),
				Region:           aws.String(config.Region),
				Endpoint:         aws.String(config.Address),
				S3ForcePathStyle: aws.Bool(true),
			},
			Profile: config.Profile,
		},
	)
}

