package queue

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/razorpay/devstack/hooks/sqs-configurator/constants"
	"go.uber.org/zap"
	"net/http"
	"strings"
	"time"
)

// Queue holds the handler for the client and auxiliary info
type Queue struct {
	client sqsiface.SQSAPI
	prefix string
}

// CreateQueue creates a queue based on the config
func (q Queue) CreateQueue(name string) (string, error) {

	if name == "" {
		zap.L().Debug("You must supply a queue name (-q QUEUE")
		return "", errors.New("Queue Name empty.")
	}

	result, err := q.client.CreateQueue(&sqs.CreateQueueInput{
		QueueName: &name,
		Attributes: map[string]*string{
			"DelaySeconds":           aws.String(constants.QUEUE_DELAY_SECONDS),
			"MessageRetentionPeriod": aws.String(constants.QUEUE_MESSAGE_RETENTION_PERIOD),
		},
	})

	if err != nil {
		return "", err
	}

	return *result.QueueUrl, nil
}

// New initializes the client with the provider
func New(Provider string) (Queue, error) {
	var q Queue
	var prefix string

	// Base config
	awsConfig := &aws.Config{
		Region:                        aws.String(constants.AWS_REGION),
		HTTPClient:                    newHTTPClientWithSettings(),
		CredentialsChainVerboseErrors: aws.Bool(true),
		MaxRetries:                    aws.Int(constants.MAX_RETRIES),
	}

	//Configuring AWS config for different providers
	if strings.EqualFold(Provider, "localstack") {
		awsConfig.Endpoint = aws.String(constants.LOCALSTACK_ENDPOINT)
		awsConfig.Credentials = credentials.AnonymousCredentials
		prefix = constants.LOCALSTACK_PREFIX
	} else if strings.EqualFold(Provider, "AWS") {
		//to-do , use operator instead of the code revisit if using this for AWS
		awsConfig.Endpoint = aws.String(constants.AWS_ENDPOINT)
		prefix = constants.AWS_PREFIX
	} else {
		return q, errors.New("Unidentified Provider")
	}

	sqsClient := sqs.New(session.Must(session.NewSession(awsConfig)))

	q = Queue{
		client: sqsClient,
		prefix: prefix,
	}

	return q, nil
}

// newHTTPClientWithSettings creates new HTTP client using constant settings
func newHTTPClientWithSettings() *http.Client {

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 0
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 100

	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	return httpClient
}
