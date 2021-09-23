package cloud

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

var _ MessageClient = SQS{}

type SQS struct {
	timeout time.Duration
	client  *sqs.SQS
}

func NewSQS(session *session.Session, timeout time.Duration) SQS {
	return SQS{
		timeout: timeout,
		client:  sqs.New(session),
	}
}

func (s SQS) CreateQueue(ctx context.Context, queueName string, isDLX bool) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	ret := "345600" // 4 days
	if isDLX {
		ret = "1209600" // 14 days
	}

	res, err := s.client.CreateQueueWithContext(ctx, &sqs.CreateQueueInput{
		QueueName: aws.String(queueName),
		Attributes: map[string]*string{
			"MessageRetentionPeriod":        aws.String(ret),
			"VisibilityTimeout":             aws.String("5"),
			"ReceiveMessageWaitTimeSeconds": aws.String("20"), // Enable long polling
		},
	})
	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}

	return *res.QueueUrl, nil
}

func (s SQS) QueueARN(ctx context.Context, queueURL string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.GetQueueAttributesWithContext(ctx, &sqs.GetQueueAttributesInput{
		AttributeNames: []*string{aws.String("QueueArn")},
		QueueUrl:       aws.String(queueURL),
	})
	if err != nil {
		return "", fmt.Errorf("get attributes: %w", err)
	}

	if len(res.Attributes) != 1 {
		return "", fmt.Errorf("not found")
	}

	return *res.Attributes["QueueArn"], nil
}

func (s SQS) BindDLX(ctx context.Context, queueURL, dlxARN string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	policy, err := json.Marshal(map[string]string{
		"deadLetterTargetArn": dlxARN,
		"maxReceiveCount":     "3",
	})
	if err != nil {
		return fmt.Errorf("marshal policy")
	}

	if _, err := s.client.SetQueueAttributesWithContext(ctx, &sqs.SetQueueAttributesInput{
		QueueUrl: aws.String(queueURL),
		Attributes: map[string]*string{
			sqs.QueueAttributeNameRedrivePolicy: aws.String(string(policy)),
		},
	}); err != nil {
		return fmt.Errorf("set attributes: %w", err)
	}

	return nil
}

func (s SQS) Send(ctx context.Context, req *SendRequest) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	attrs := make(map[string]*sqs.MessageAttributeValue, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs[attr.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	res, err := s.client.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageAttributes: attrs,
		MessageBody:       aws.String(req.Body),
		QueueUrl:          aws.String(req.QueueURL),
	})
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	return *res.MessageId, nil
}

func (s SQS) Receive(ctx context.Context, queueURL string) (*Message, error) {
	// timeout = WaitTimeSeconds + 5
	ctx, cancel := context.WithTimeout(ctx, time.Second*(20+5))
	defer cancel()

	res, err := s.client.ReceiveMessageWithContext(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:              aws.String(queueURL),
		MaxNumberOfMessages:   aws.Int64(1),
		WaitTimeSeconds:       aws.Int64(20),
		MessageAttributeNames: aws.StringSlice([]string{"All"}),
	})
	if err != nil {
		return nil, fmt.Errorf("receive: %w", err)
	}

	if len(res.Messages) == 0 {
		return nil, nil
	}

	attrs := make(map[string]string)
	for key, attr := range res.Messages[0].MessageAttributes {
		attrs[key] = *attr.StringValue
	}

	return &Message{
		ID:            *res.Messages[0].MessageId,
		ReceiptHandle: *res.Messages[0].ReceiptHandle,
		Body:          *res.Messages[0].Body,
		Attributes:    attrs,
	}, nil
}

func (s SQS) Delete(ctx context.Context, queueURL, rcvHandle string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.DeleteMessageWithContext(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueURL),
		ReceiptHandle: aws.String(rcvHandle),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}
