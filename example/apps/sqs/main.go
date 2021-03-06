package main

import (
	"log"
	"time"
    "github.com/razorpay/devstack/example/apps/sqs/internal/message"
	"github.com/razorpay/devstack/example/apps/sqs/internal/pkg/cloud"
)

func main() {
	// Create a session instance.
	awsConfig := cloud.NewConfig()
	sqs, err := cloud.New(cloud.Config{
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
		message.Message(cloud.NewSQS(sqs, time.Second*5))
		time.Sleep(time.Second)
		log.Printf("Resending another message")
	}
}

