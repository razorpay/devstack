package cloud

import (
"context"
)

type MessageClient interface {
	// Creates a new long polling queue and returns its URL.
	CreateQueue(ctx context.Context, queueName string, isDLX bool) (string, error)
	// Get a queue ARN.
	QueueARN(ctx context.Context, queueURL string) (string, error)
	// Binds a DLX queue to a normal queue.
	BindDLX(ctx context.Context, queueURL, dlxARN string) error
	// Send a message to queue and returns its message ID.
	Send(ctx context.Context, req *SendRequest) (string, error)
	// Long polls given amount of messages from a queue.
	Receive(ctx context.Context, queueURL string) (*Message, error)
	// Deletes a message from a queue.
	Delete(ctx context.Context, queueURL, rcvHandle string) error
}
