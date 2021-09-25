// Package config will manage all application level configurations
// config file will be taken based on the application environment
// all the configuration available in governor file will be overwritten
// this will be immutable as it always provides the value of the struct
package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/tracing"
	"github.com/sirupsen/logrus"
)

const (
	// FilePath - relative path to the config directory
	FilePath = "%s/config/%s"

	// DefaultFilename - Filename format of default config file
	DefaultFilename = "conf.default.toml"

	// EnvFilename - Filename format of env specific config file
	EnvFilename = "conf.%s.toml"
)

var (
	// config : this will hold all the application configuration
	config      AppConfig
	appBasePath string
)

// appConfig global configuration struct definition
type AppConfig struct {
	Tracing tracing.Config `toml:"tracing"`
}

// LoadConfig will load the configuration available in the cnf directory available in basePath
// conf file will takes based on the env provided
// governor file content will override the conf available
func LoadConfig(basePath string, env string) {

	appBasePath = basePath

	// reading conf based on default environment
	loadConfigFromFile(basePath, DefaultFilename, "")

	// reading env file and override conf values; if env file exists
	loadConfigFromFile(basePath, EnvFilename, env)

	//Validate
	ValidateConfig("", GetConfig())
}

func ValidateConfig(field string, a interface{}) {
	v := reflect.ValueOf(a)
	for j := 0; j < v.NumField(); j++ {
		f := v.Field(j)
		n := v.Type().Field(j).Name
		switch f.Kind() {
		case reflect.String:
			if f.String() == "" {
				logrus.WithError(errors.New("ENV VARIABLE not found for " + field + n)).Panic("InvalidConfigType")
				return
			}
		case reflect.Struct:
			ValidateConfig(field+n+".", f.Interface())
		}
	}
}

// GetConfig : will give the struct as value so that the actual conf doesn't get tampered
func GetConfig() AppConfig {
	return config
}

// loadConfigFromFile: load/overwrite config values from given basepath, filename and env
func loadConfigFromFile(basePath string, filename string, env string) {
	path := getFilePath(basePath, filename, env)
	_, err := os.Stat(path)
	if err != nil {
		logrus.WithError(err).Error("ConfigFileNotFound")
	}

	content := readConfigFile(path)

	// inject environment variables
	content = os.ExpandEnv(string(content))

	if _, err := toml.Decode(content, &config); err != nil {
		logrus.WithError(err).Panic("InvalidConfigType")
	}
}

// readConfigFile: config file will be read from the given file and gives the content in string format
func readConfigFile(path string) string {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		logrus.WithError(err).Error("ConfigFileNotFound")
		return ""
	}

	return string(data)
}

// getFilePath: gives the file path based on the environment provided
// file path will be relative to the application and determined by basePath
func getFilePath(basePath string, fileName string, env string) string {
	if env != "" {
		fileName = fmt.Sprintf(fileName, env)
	}

	path := fmt.Sprintf(FilePath, basePath, fileName)

	return path
}
