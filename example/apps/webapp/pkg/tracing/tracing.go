package tracing

import (
	"fmt"
	"io"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics"
)

// Tracer instance
type JaegerTracer struct {
	// Tracer is a simple, thin interface for Span creation and SpanContext
	Tracer opentracing.Tracer
	// Interface for wrapping close methods when closing the tracer
	Closer io.Closer
}

var (
	config Config

	tracers []*JaegerTracer

	DBPeerService string
	// Database HostName
	DBHostName string
	// Database port
	DBPort uint16

	RedisPeerService string
	// Redis Hostname
	RedisHostName string
	// Redis Port
	RedisPort uint16

	// Don't remove the following. This will be used for deployment markers. RTFM
	ServiceVersion string //Git Commit Hash / Service Version

	ContextPropogationKey          string
	ContextKey                     string
	DefaultContextPropogationValue string
	BaggagePrefix                  string
	BaggageKeys                    string
)

var (
	tracer = &JaegerTracer{
		Tracer: nil,
		Closer: nil,
	}
)

// Exclude the following request headers from going into the spans from the http headers
var RequestHeaderExclusions = map[string]bool{}

//Exclude the following query params from going into the spans from the http headers
var RequestQueryParamExclusions = map[string]bool{}

var RouteURLsToExclude []string

//Gets the tracer object. Tracer is a singleton
func GetTracer() *JaegerTracer {
	return tracer
}

//Return Closer interface. Useful for flushing the spans when needed or closing it during application shutdown
func (tr *JaegerTracer) GetCloser() io.Closer {
	return tr.Closer
}

//Get the opentracing object from the tracer
func (tr *JaegerTracer) GetTracer() opentracing.Tracer {
	return tr.Tracer
}

func InitConfig(conf Config) {
	config = conf
}

/*
 * General Tracer Initializer.
 * To be called from init of any datastore like mysql/redis etc
 *
 * Environment determines the actual appearance on the jaeger UI
 * (e.g. {application}-prod, {application}-stage, {application}-canary etc
 */
func InitTracing(serviceName string) *JaegerTracer {
	jaegerEnabled := config.Enabled
	jaegerHost := config.JaegerHostName
	jaegerPort := config.JaegerPort
	debugMode := config.EnableDebug
	RequestHeaderExclusions = config.RequestHeaderExclusions
	RouteURLsToExclude = config.RouteURLsToExclude
	RequestQueryParamExclusions = config.RequestQueryParamExclusions

	logSpans := false
	if debugMode == true {
		logSpans = true
	}
	serviceName = fmt.Sprintf("%s-%s", serviceName, config.Env)
	if config.ServiceVersion == "" {
		ServiceVersion = DefaultServiceVersion
	} else {
		ServiceVersion = config.ServiceVersion
	}

	jaegerHostPort := fmt.Sprintf("%s:%v", jaegerHost, jaegerPort)
	cfg := &jaegercfg.Configuration{
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},

		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:           logSpans,
			LocalAgentHostPort: jaegerHostPort,
		},

		ServiceName: serviceName,
		Disabled:    !jaegerEnabled,
		Tags: []opentracing.Tag{
			{ServiceVersionKey, ServiceVersion},
			{Environment, config.Env},
		},
	}

	jLogger := config.Logger

	// set your own metrics factory here
	jMetricsFactory := metrics.NullFactory

	// Initialize tracer with a logger and a metrics factory
	tr, cl, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		logrus.Fatalf("Error:%v", err)
	}

	jTracer := &JaegerTracer{Tracer: tr, Closer: cl}
	// keep track of all inited tracers to be able to close when required.
	tracers = append(tracers, jTracer)
	return jTracer
}

/*
 * Initializes global app tracer Initializer and sets tracing configs of peer services also.
 *
 * To be called from main/init of the application, before calling InitTracing from anywhere else.
 * opentracing.StartSpan would start span on this tracer, by default
 */
func InitAppTracer(serviceName string) *JaegerTracer {
	tracer = InitTracing(serviceName)

	assignDbTracingConfig(config)
	assignRedisTracingConfig(config)
	assignContextConfig(config.Context)

	opentracing.SetGlobalTracer(tracer.Tracer)
	return tracer
}

// close all tracers which are opened in this app.
func CloseTracers() {
	for _, tr := range tracers {
		tr.Closer.Close()
	}
}

func assignContextConfig(contextConfig Context) {
	ContextPropogationKey = contextConfig.ContextPropogationKey
	ContextKey = contextConfig.ContextKey
	DefaultContextPropogationValue = contextConfig.DefaultContextPropogationValue
	BaggagePrefix = contextConfig.BaggagePrefix
	BaggageKeys = contextConfig.BaggageKeys
}

func assignRedisTracingConfig(conf Config) {
	RedisPeerService = conf.Redis.PeerService
	RedisHostName = conf.Redis.HostName
	RedisPort = conf.Redis.Port
}

func assignDbTracingConfig(conf Config) {
	DBPeerService = conf.Database.PeerService
	DBHostName = conf.Database.HostName
	DBPort = conf.Database.Port
}

// Helper method for adding a map of client tags to the given span object
func AddClientTags(clientSpan opentracing.Span, tags map[string]interface{}) opentracing.Span {
	if clientSpan != nil {
		span := clientSpan
		for key, value := range tags {
			span = span.SetTag(key, value)
		}
		return span
	}
	return clientSpan
}

// Helper method for injecting traces for remote http calls
func InjectClientTrace(clientSpan opentracing.Span, req *http.Request) {
	tracer.Tracer.Inject(clientSpan.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))
}

// Another helper method for injecting traces explictly with request object, url and http method
func InjectClientTraceWithUrlMethod(clientSpan opentracing.Span, req *http.Request, url string, method string) {
	ext.SpanKindRPCClient.Set(clientSpan)
	ext.HTTPUrl.Set(clientSpan, url)
	ext.HTTPMethod.Set(clientSpan, method)
	InjectClientTrace(clientSpan, req)
}

// Helper method to record application errors into traces
func RecordError(clientSpan opentracing.Span, err error, clientName string) {
	if clientSpan != nil {
		ext.Error.Set(clientSpan, true)
		clientSpan.LogKV("error", err)
	}
}
