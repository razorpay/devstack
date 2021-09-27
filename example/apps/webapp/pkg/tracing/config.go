package tracing

import "github.com/prometheus/client_golang/prometheus"

type Config struct {
	ServiceName                 string
	ServiceVersion              string
	Env                         string
	JaegerHostName              string
	JaegerPort                  string
	Enabled                     bool
	EnableDebug                 bool
	RequestHeaderExclusions     map[string]bool
	RequestQueryParamExclusions map[string]bool
	RouteURLsToExclude          []string
	Logger                      Logger
	PrometheusRegister          prometheus.Registerer
	Database                    DbConfig
	Redis                       RedisConfig
	Context                     Context
}

type DbConfig struct {
	PeerService string
	HostName    string
	Port        uint16
}

type RedisConfig struct {
	PeerService string
	HostName    string
	Port        uint16
}

type Context struct {
	ContextPropogationKey          string
	ContextKey                     string
	DefaultContextPropogationValue string
	BaggagePrefix                  string
	BaggageKeys                    string
}
