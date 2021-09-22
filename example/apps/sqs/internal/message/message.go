package message

import (
	"context"
	"log"
	"time"
	"fmt"

	"github.com/razorpay/devstack/example/apps/sqs/internal/pkg/cloud"
)

func Message(client cloud.MessageClient) {
	ctx := context.Background()

	dlxURL := createQueueDLX(ctx, client)
	queURL := createQueue(ctx, client)
	dlxARN := queueARN(ctx, client, dlxURL)
	bindDLX(ctx, client, queURL, dlxARN)
	send(ctx, client, queURL)
	rcvHnd := receive(ctx, client, queURL)
	deleteMessage(ctx, client, queURL, rcvHnd)
}

func createQueueDLX(ctx context.Context, client cloud.MessageClient) string {
	url, err := client.CreateQueue(ctx, "welcome-email-queue.dlx", true)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("create queue:", url)

	return url
}

func createQueue(ctx context.Context, client cloud.MessageClient) string {
	url, err := client.CreateQueue(ctx, "welcome-email-queue", false)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("create queue:", url)

	return url
}

func queueARN(ctx context.Context, client cloud.MessageClient, url string) string {
	arn, err := client.QueueARN(ctx, url)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("queue ARN:", arn)

	return arn
}

func bindDLX(ctx context.Context, client cloud.MessageClient, queueURL, dlxARN string) {
	if err := client.BindDLX(ctx, queueURL, dlxARN); err != nil {
		log.Fatalln(err)
	}
	log.Println("bind DLX: ok")
}

func send(ctx context.Context, client cloud.MessageClient, queueURL string) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	id, err := client.Send(ctx, &cloud.SendRequest{
		QueueURL: queueURL,
		Body:     "Message body!",
		Attributes: []cloud.Attribute{
			{
				Key:   "Title",
				Value: "SQS send message",
				Type:  "String",
			},
			{
				Key:   "Timestamp",
				Value: fmt.Sprintf("%s", timestamp),
				Type:  "Number",
			},
		},
	})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("send: message ID:", id)
}

func receive(ctx context.Context, client cloud.MessageClient, queueURL string) string {
	res, err := client.Receive(ctx, queueURL)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("receive:", res)

	return res.ReceiptHandle
}

func deleteMessage(ctx context.Context, client cloud.MessageClient, queueURL, rcvHnd string) {
	if err := client.Delete(ctx, queueURL, rcvHnd); err != nil {
		log.Fatalln(err)
	}
	log.Println("delete message: ok")
}

