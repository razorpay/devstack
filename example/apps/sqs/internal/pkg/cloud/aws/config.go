package aws

import "os"

//SQSConfig struct holds the basic config details of SQS Queue
type SQSConfig struct {
	Address string
	Region string
	Profile string
	AwsKey string
	AwsSecret string
}

//New config Constructor
func NewConfig() *SQSConfig {
	return &SQSConfig{
		Address: getEnv("SQS_URL", "http://localhost:4566"),
		Region: getEnv("AWS_REGION", "ap-south-1"),
		Profile: getEnv("AWS_PROFILE", "localstack"),
		AwsKey: getEnv("AWS_KEY", "test"),
		AwsSecret: getEnv("AWS_SECRET", "test"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

