package helpers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	tracing "github.com/razorpay/devstack/example/apps/webapp/pkg/tracing"
)

//Creates a new operational child span from the gin context with attributes/tagging. Use this to create spans from gin controllers
func GetChildSpanFromGinContext(c *gin.Context, operation string, tags map[string]interface{}) (opentracing.Span, context.Context) {
	ctxInterface, exists := c.Get(tracing.ContextKey)
	if ctxInterface == nil || !exists {
		// This means the middleware didn't kick in. In all these cases, just initialize a noop tracer and send a span out
		tr := opentracing.NoopTracer{}
		return opentracing.StartSpanFromContextWithTracer(c.Request.Context(), tr, operation)
	}
	ctx := ctxInterface.(context.Context)
	childSpan, ctx := GetChildSpanFromContext(ctx, operation, tags)
	c.Set(tracing.ContextKey, ctx)
	return childSpan, ctx
}

//Creates a new operational child span from the golang context with attributes/tagging. Use this to create spans from existing go context object
func GetChildSpanFromContext(ctx context.Context, operation string, tags map[string]interface{}) (opentracing.Span, context.Context) {
	parentSpan := opentracing.SpanFromContext(ctx)
	if parentSpan == nil {
		return nil, ctx
	}
	childSpan := opentracing.StartSpan(operation, opentracing.ChildOf(parentSpan.Context()))
	if len(tags) > 0 {
		childSpan = tracing.AddClientTags(childSpan, tags)
	}

	ctxPropogationValue := ctx.Value(tracing.ContextPropogationKey)
	if ctxPropogationValue == nil {
		ctxPropogationValue = tracing.DefaultContextPropogationValue
	}
	tmp := ctx.Value(tracing.BaggageKeys)
	var baggageKeys []string
	if tmp != nil {
		baggageKeys = tmp.([]string)
		for _, bKey := range baggageKeys {
			bValue := parentSpan.BaggageItem(bKey)
			childSpan.SetBaggageItem(bKey, bValue)
		}
	}

	// Now, propogate the baggage keys across the context
	ctx = context.WithValue(ctx, tracing.BaggageKeys, baggageKeys)
	ctx = context.WithValue(ctx, tracing.ContextPropogationKey, ctxPropogationValue)
	childSpan.SetTag(tracing.ContextPropogationKey, ctxPropogationValue)

	return childSpan, opentracing.ContextWithSpan(ctx, childSpan)
}

//Creates a new operational child span from the golang context, specifically for redis operations
func StartRedisSpanFromContext(ctx context.Context, operation string, tags map[string]interface{}) (opentracing.Span, context.Context) {
	childSpan, ctx := GetChildSpanFromContext(ctx, operation, tags)
	if childSpan == nil {
		return nil, ctx
	}
	ext.SpanKindRPCClient.Set(childSpan)
	ext.PeerService.Set(childSpan, tracing.RedisPeerService)
	ext.PeerAddress.Set(childSpan, tracing.RedisHostName)
	ext.PeerPort.Set(childSpan, tracing.RedisPort)
	ext.DBType.Set(childSpan, "redis")
	return childSpan, opentracing.ContextWithSpan(ctx, childSpan)
}
