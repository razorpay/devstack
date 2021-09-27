package main

import (
	"flag"
	"github.com/razorpay/devstack/hooks/secret_cloner/controllers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

// Params defines the configuration for the application actions
type Params struct {
	NameFlag      string
	DebugFlag     bool
	IsLocal       bool
}

//SecretConfig defines the details of secret to be acted upon
type SecretConfig struct {
	Action               string
	Namespace            string
	SecretName           string
	SecretSuffix         string
	UpdateEntries        map[string]Secret
	Ttl                  string
}
//Secret stores the secret details
type Secret struct {
	Key string
	Value string
}

var (
 	params     Params
 	secretData SecretConfig
 	failure    bool
 	logger      *zap.Logger
 	err        error
)


func main(){
	params.DebugFlag = *flag.Bool("debug",false,"To log all events")
	params.IsLocal = *flag.Bool("local",false,"Local debuggig")

	flag.Parse()
	initialize()
	initializeConfig()
	if strings.EqualFold(secretData.Action,"clone"){
		processClone()
	} else if strings.EqualFold(secretData.Action,"update"){
		processUpdate()
	} else {
		logger.Info("Unsupported Action")
	}
}

//initializeConfig initializes the config using viper
func initializeConfig() {
	//load the config  from env in case of cluster execution
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.Unmarshal(&secretData)
}

//initialize the application
func initialize(){
	var cfg zap.Config
	if params.DebugFlag {
		cfg = getConfig(zapcore.DebugLevel)
	} else {
		cfg = getConfig(zapcore.InfoLevel)
	}
	logger, err = cfg.Build()

	if err != nil {
		panic(err)
	}
	defer logger.Sync()

}

//getConfig internal method for reading zap configuration
func getConfig(level zapcore.Level) zap.Config {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(level),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:   zapcore.EncoderConfig{
			MessageKey: "message"},
	}
	return cfg
}

//processClone clones the secret , deletes if the secret already exists
func processClone(){
	err := controllers.New(params.IsLocal)
	if err != nil{
		logger.Info("Initialziation error. Error : "+ err.Error())
		os.Exit(1)
	}
    oldSecret , err := controllers.GetSecret(secretData.Namespace,secretData.SecretName)
    if err != nil {
    	logger.Info("The secret with the name not found exiting " + secretData.SecretName)
    	logger.Debug(err.Error())
    	os.Exit(1)
	}
	annotations := map[string]string{"janitor/ttl":secretData.Ttl}
	metaObject := metaV1.ObjectMeta{Name: secretData.SecretName+"-"+secretData.SecretSuffix , Annotations: annotations}
    newSecret := v1.Secret{Data: oldSecret.Data,ObjectMeta: metaObject}
    secret ,_ := controllers.GetSecret(secretData.Namespace,newSecret.GetName())
    if  secret != nil {
		err := controllers.DeleteSecret(secretData.Namespace, newSecret.GetName())
		if err != nil {
			logger.Info("The deletion of secret failed Error: "+ err.Error())
			os.Exit(1)
		}
	}
	err = controllers.CreateSecret(secretData.Namespace, &newSecret)
	if err != nil {
		logger.Info("The creation of secert failed Error: "+err.Error())
		os.Exit(1)
	}
}

//processUpdate updates the secrets with the provided values
func processUpdate(){
	failure = false

	for _,value := range secretData.UpdateEntries {
		err := controllers.UpdateSecret(secretData.Namespace,secretData.SecretName,value.Key,value.Value)
		if err != nil {
			logger.Info("The creation of secert failed Error: "+err.Error())
			failure=true
		}
	}

	if failure {
		logger.Info("Failing as some of the components creation failed with error , please check the above logs")
		os.Exit(1)
	}
}
