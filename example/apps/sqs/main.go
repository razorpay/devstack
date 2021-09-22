package main

import (
	"log"
	"time"
    "github.com/razorpay/devstack/example/apps/sqs/internal/message"
	"github.com/razorpay/devstack/example/apps/sqs/internal/pkg/cloud/aws"
)

func main() {
	// Create a session instance.
	awsConfig := aws.NewConfig()
	sqs, err := aws.New(aws.Config{
		Address: awsConfig.Address,
		Region: awsConfig.Region,
		Profile: awsConfig.Profile,
		ID: awsConfig.AwsKey,
		Secret: awsConfig.AwsSecret,
	})

	if err != nil {
		log.Fatalln(err)
	}

	// Test message
	for {
		message.Message(aws.NewSQS(sqs, time.Second*5))
		time.Sleep(time.Second)
		log.Printf("Resending another message")
	}
}

