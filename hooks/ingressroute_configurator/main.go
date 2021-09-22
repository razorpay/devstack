package main

import (
	"flag"
	"github.com/razorpay/devstack/hooks/ingressroute_configurator/traefik"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"path/filepath"
	"strings"
)


// Params defines the configuration for the application actions
type Params struct {
	NameFlag *string
	IsLocal  *bool
	Debug    *bool
	Action   *string
}

// IngressRouteConfig defines the details of the ingress route object being worked upon
type IngressRouteConfig struct {
	IngressUrl string
	Namespace  string
	HeaderValue string
	IngressRouteName string
	ServiceName string
	ServicePort int32
}


var (
	params Params

	irConfig IngressRouteConfig

	clientSet *kubernetes.Clientset

	err error

	log *zap.Logger
)


func main(){
	params.Action  = flag.String("action","","The action that needs to be done , needs to either of Update/Delete")
	params.IsLocal = flag.Bool("local",false,"Is the config run in local Default : false")
	params.Debug   = flag.Bool("debug",false,"To log all the events Default : false")
	flag.Parse()
	Initialize()
	InitializeConfig()
	InitializeClient()
   if strings.EqualFold(*params.Action,"update"){
		ProcessUpdate()
	} else if strings.EqualFold(*params.Action,"delete"){
		ProcessDelete()
	} else {
		log.Info("Unsupported Action")
	}

}

//Initialize the application
func Initialize(){
	var cfg zap.Config
	if *params.Debug {
		cfg = getConfig(zapcore.DebugLevel)
	} else {
		cfg = getConfig(zapcore.InfoLevel)
	}
	log, err = cfg.Build()

	if err != nil {
		panic(err)
	}
	defer log.Sync()

	if *params.Action== "" {
		log.Info("Invalid Parameters\n")
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
func InitializeConfig() {
	//load the config  from env in case of cluster execution
	viper.SetConfigName("app")
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.ReadInConfig()
	viper.Unmarshal(&irConfig)
}

//InitializeClient initializes the K8s client
func InitializeClient() {
	var config *rest.Config

	if *params.IsLocal {
		home, exists := os.LookupEnv("HOME")
		if !exists {
			home = "/root"
		}
		configPath := filepath.Join(home, ".kube", "config")
		config, _ = clientcmd.BuildConfigFromFlags("", configPath)
	} else  {
		config, _ = rest.InClusterConfig()
	}

	clientSet, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Panic("Failed to create K8s clientset")
	}

}

//ProcessUpdate processes the given request , follows a looping approach to avoid concurrency overriding
func ProcessUpdate(){
	retry := true
	retryCount := 0
	for retry && retryCount < traefik.RETRY_LIMIT {
		ingressRoute, _ := traefik.GetIngressRoute(clientSet, irConfig.Namespace, irConfig.IngressRouteName)
		if traefik.IsRulePresent(irConfig.IngressUrl, irConfig.HeaderValue, ingressRoute.Spec) {
			log.Info("The header rule for the key is already present")
			retry = false
		} else {
			traefik.InsertMatchRule(irConfig.IngressUrl, irConfig.HeaderValue, irConfig.ServiceName, irConfig.ServicePort , &ingressRoute.Spec)
			_, err = traefik.PatchIngressRoute(clientSet, irConfig.Namespace, irConfig.IngressRouteName, *ingressRoute)
			if err != nil{
				log.Info("Error in updating the ingress route , retrying")
				log.Info(err.Error())
			}
			retryCount++
		}
	}
	if retryCount == traefik.RETRY_LIMIT {
		log.Info("Failing as some of the components creation failed with error , please check the above logs")
		os.Exit(1)
	}
}

//ProcessDelete processes the delete request , follows a looping approach to avoid concurrency overriding
func ProcessDelete(){
	retry := true
	retryCount := 0
	for retry && retryCount < traefik.RETRY_LIMIT {
		ingressRoute, _ := traefik.GetIngressRoute(clientSet, irConfig.Namespace, irConfig.IngressRouteName)
		if traefik.IsRulePresent(irConfig.IngressUrl, irConfig.HeaderValue, ingressRoute.Spec) {
			err := traefik.DeleteMatchRule(irConfig.IngressUrl, irConfig.HeaderValue, &ingressRoute.Spec)
			if err != nil {
				log.Info("The header rule for the key is not present exiting ")
			}
			_, err = traefik.PatchIngressRoute(clientSet, irConfig.Namespace, irConfig.IngressRouteName, *ingressRoute)
			if err != nil {
				log.Info("Error in updating the ingress route , retrying")
			}
			retryCount++
		} else {
			log.Info("The header rule for the key is not present exiting")
			retry = false
		}
	}
	if retryCount == traefik.RETRY_LIMIT {
		log.Info("Failing as some of the components creation failed with error , please check the above logs")
		os.Exit(1)
	}
}
