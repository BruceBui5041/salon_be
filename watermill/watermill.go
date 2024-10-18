package watermill

import (
	"context"
	"video_server/component"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/pubsub/gochannel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type MessageHandler func(appCtx component.AppContext, msg *message.Message)

type localPubSub struct {
	blockPubSub   *gochannel.GoChannel
	unblockPubSub *gochannel.GoChannel
}

func NewPubsubPublisher() *localPubSub {
	return &localPubSub{
		blockPubSub: gochannel.NewGoChannel(
			gochannel.Config{
				BlockPublishUntilSubscriberAck: true,
			},
			watermill.NewStdLogger(false, false),
		),
		unblockPubSub: gochannel.NewGoChannel(
			gochannel.Config{
				BlockPublishUntilSubscriberAck: false,
			},
			watermill.NewStdLogger(false, false),
		),
	}
}

func (ps *localPubSub) GetUnblockPubSub() *gochannel.GoChannel {
	return ps.unblockPubSub
}

func (ps *localPubSub) GetBlockPubSub() *gochannel.GoChannel {
	return ps.blockPubSub
}

// Helper function to get or generate a correlation ID
func getCorrelationID(ctx context.Context) string {
	if correlationID, ok := ctx.Value("correlation_id").(string); ok {
		return correlationID
	}
	return uuid.New().String()
}

func GetTraceAndSpanID(ctx context.Context) (string, string) {
	spanCtx := trace.SpanContextFromContext(ctx)
	return spanCtx.TraceID().String(), spanCtx.SpanID().String()
}

func setTracingMetadata(ctx context.Context, msg *message.Message) {
	traceID, spanID := GetTraceAndSpanID(ctx)
	msg.Metadata.Set("trace_id", traceID)
	msg.Metadata.Set("span_id", spanID)
	msg.Metadata.Set("correlation_id", getCorrelationID(ctx))
}

// New function to extract tracing metadata
func extractTracingMetadata(msg *message.Message) (trace.TraceID, trace.SpanID, error) {
	traceID, err := trace.TraceIDFromHex(msg.Metadata.Get("trace_id"))
	if err != nil {
		return trace.TraceID{}, trace.SpanID{}, err
	}
	spanID, err := trace.SpanIDFromHex(msg.Metadata.Get("span_id"))
	if err != nil {
		return trace.TraceID{}, trace.SpanID{}, err
	}
	return traceID, spanID, nil
}

func createTracedHandler(msg *message.Message, handlerName string) (context.Context, trace.Span) {
	ctx := msg.Context()
	tracer := otel.Tracer("pubsub")

	traceID, spanID, err := extractTracingMetadata(msg)
	if err != nil {
		// If we can't extract tracing metadata, we'll start a new trace
		ctx, span := tracer.Start(ctx, handlerName)
		return ctx, span
	}

	// Create a new span context
	spanCtx := trace.NewSpanContext(trace.SpanContextConfig{
		TraceID: traceID,
		SpanID:  spanID,
		Remote:  true,
	})

	// Start a new span as a child of the message's span
	ctx, span := tracer.Start(trace.ContextWithSpanContext(ctx, spanCtx), handlerName)
	return ctx, span
}
