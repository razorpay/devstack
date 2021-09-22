package main

import (
	"flag"
	"github.com/razorpay/devstack/hooks/sqs-configurator/package/queue"
	"github.com/razorpay/devstack/hooks/sqs-configurator/package/secrets"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

// Queue defines the queue details and the mapping secret
type Queue struct {
	Name      string `yaml:name`
	SecretKey string `yaml:secret-key`
}

//QueueConfig defines the config for the configurator
type QueueConfig struct {
	Queue                map[string]Queue
	UpdateSecret         bool
	KubeSecret           string
	Namespace            string
	Provider             string
	EnableEndpointPrefix bool
}

// Params defines the configuration for the application actions
type Params struct {
	NameFlag  *string
	DebugFlag *bool
	IsLocal   *bool
}

var (
	paramsFlag Params
	queueData  QueueConfig
	log        *zap.Logger
	err        error
)

func main() {
	paramsFlag.NameFlag = flag.String("name", "test", "Config File Name")
	paramsFlag.DebugFlag = flag.Bool("debug", false, "To log all events")
	paramsFlag.IsLocal = flag.Bool("local", false, "Local debugging")

	flag.Parse()
	initialize()
	initializeConfig()
	Process()
}

//initialize the application
func initialize() {
	var cfg zap.Config
	if *paramsFlag.DebugFlag {
		cfg = getConfig(zapcore.DebugLevel)
	} else {
		cfg = getConfig(zapcore.InfoLevel)
	}
	log, err = cfg.Build()

	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if *paramsFlag.NameFlag == "" {
		log.Info("inv", zap.String("IP", "Invalid parameters provided"))
		flag.PrintDefaults()
		os.Exit(1)
	}
}

//getConfig internal method for reading zap configuration
func getConfig(level zapcore.Level) zap.Config {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}
	return cfg
}

//InitializeConfig initializes the config using viper
func initializeConfig() {
	//load the config  from env in case of cluster execution
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.Unmarshal(&queueData)
	log.Debug("QueueData", zap.Any("qd", queueData))
}

// Process the sqs configurator request
func Process() {
	sqs, err := queue.New(queueData.Provider)
	if err != nil {
		log.Debug("failure", zap.String("FAILURE", "Error in creating SQS client"))
		os.Exit(1)
	}

	//iterate over all the queues data provided
	for _, value := range queueData.Queue {
		qUrl, err := sqs.CreateQueue(value.Name)
		if err != nil {
			log.Debug("failure", zap.String("FAILURE", "Failed creating SQS queue for queue with name "+value.Name+err.Error()))
			os.Exit(1)
		}
		//update the secret only if update secret is true
		if queueData.UpdateSecret {
			var queueVal string
			//Check if the endpoint needs to be prefixed
			if queueData.EnableEndpointPrefix {
				queueVal = qUrl
			} else {
				queueVal = value.Name
			}
			err = secrets.UpdateSecret(queueData.Namespace, queueData.KubeSecret, value.SecretKey, queueVal, paramsFlag.IsLocal)
			if err != nil {
				log.Debug("failure", zap.String("FAILURE", "failed to create key "+err.Error()))
				os.Exit(1)
			}
		}
	}
	return
}
