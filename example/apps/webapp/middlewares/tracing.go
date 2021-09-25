package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/razorpay/devstack/example/apps/webapp/pkg/tracing"
)

func TracingMiddleware() gin.HandlerFunc {
	var serverSpan opentracing.Span
	jTracer := tracing.GetTracer()
	return func(c *gin.Context) {
		spanName := c.FullPath()

		// exclude urls from tracing, if any are configured
		for _, excludeUrl := range tracing.RouteURLsToExclude {
			if spanName == excludeUrl {
				// execute other middlewares and return without creating span
				c.Next()
				return
			}
		}

		spanCtx, _ := jTracer.Tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request.Header))
		serverSpan = jTracer.Tracer.StartSpan(spanName, ext.RPCServerOption(spanCtx))
		serverSpan.SetTag(tracing.HttpPath, spanName)
		serverSpan.SetTag(tracing.HttpClientHost, c.Request.Host)
		serverSpan.SetTag(tracing.HttpRequestMethod, c.Request.Method)
		// Set all headers as attributes
		var baggageKeys []string
		for key, values := range c.Request.Header {
			if _, ok := tracing.RequestHeaderExclusions[key]; !ok {
				hKey := fmt.Sprintf("%v.%v", tracing.HttpRequestHeaderPrefix, key)
				value := strings.Join(values, ",")
				serverSpan = serverSpan.SetTag(hKey, value)
				if strings.HasPrefix(key, tracing.BaggagePrefix) {
					bKey := strings.Trim(key, tracing.BaggagePrefix)
					baggageKeys = append(baggageKeys, bKey)
					serverSpan.SetBaggageItem(bKey, value)
				}
			}
		}
		//Set all query parameters as tags
		for key, values := range c.Request.URL.Query() {
			if _, ok := tracing.RequestQueryParamExclusions[key]; !ok {
				hKey := fmt.Sprintf("%v.%v", tracing.HttpRequestQueryParamPrefix, key)
				value := strings.Join(values, ",")
				serverSpan = serverSpan.SetTag(hKey, value)
			}
		}
		// Todo: Add propogation Spans
		//ctx := context.WithValue(c.Request.Context(), tracing.ParentSpanKey, serverSpan)
		ctx := opentracing.ContextWithSpan(c.Request.Context(), serverSpan)
		//Context Propogation Key
		ctxPropogationValue := c.Request.Header.Get(tracing.ContextPropogationKey)
		if ctxPropogationValue == "" {
			ctxPropogationValue = tracing.DefaultContextPropogationValue
			hKey := fmt.Sprintf("%v.%v", tracing.HttpRequestHeaderPrefix, tracing.ContextPropogationKey)
			serverSpan.SetTag(hKey, ctxPropogationValue)
		}
		// Set baggage keys if present
		ctx = context.WithValue(ctx, tracing.BaggageKeys, baggageKeys)
		ctx = context.WithValue(ctx, tracing.ContextPropogationKey, ctxPropogationValue)
		c.Set("ctx", ctx)
		defer serverSpan.Finish()
		c.Request = c.Request.WithContext(ctx)
		// Process other middlewares and request
		c.Next()
		// Record the status of the call here
		statusCode := c.Writer.Status()
		serverSpan.SetTag(tracing.HttpStatusCode, strconv.Itoa(statusCode))
		if statusCode >= http.StatusInternalServerError {
			serverSpan.SetTag(tracing.Error, fmt.Errorf("%d: %s", statusCode, http.StatusText(statusCode)))
		}
		if len(c.Errors) > 0 {
			serverSpan.SetTag(tracing.Error, c.Errors.String())
		}
	}
}
