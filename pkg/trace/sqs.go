package trace

import (
	"context"

	"github.com/aws/aws-sdk-go/service/sqs"
	mimosasqs "github.com/ca-risken/common/pkg/sqs"
	"go.opentelemetry.io/contrib"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	oteltrace "go.opentelemetry.io/otel/trace"
)

func ProcessTracingHandler(serverName string, next mimosasqs.Handler) mimosasqs.Handler {
	return mimosasqs.HandlerFunc(func(ctx context.Context, msg *sqs.Message) error {
		// TODO inherit tracing from sqs
		// The values set are referred from https://opentelemetry.io/docs/reference/specification/trace/semantic_conventions/messaging/
		tracer := otel.GetTracerProvider().Tracer(tracerName,
			oteltrace.WithInstrumentationVersion(contrib.SemVersion()))
		opts := []oteltrace.SpanStartOption{
			oteltrace.WithSpanKind(oteltrace.SpanKindServer),
		}
		ctx, span := tracer.Start(ctx, serverName+" process", opts...)
		defer span.End()

		// required attributes are only set
		attributes := []attribute.KeyValue{
			attribute.String("SpanKind", "CONSUMER"),
			attribute.String("messaging.system", "AmazonSQS"),
			attribute.String("messaging.destination", serverName),
			attribute.String("messaging.destination_kind", "queue"),
			attribute.Bool("messaging.temp_destination", false),
		}
		span.SetAttributes(attributes...)
		err := next.HandleMessage(ctx, msg)
		if err != nil {
			span.SetStatus(codes.Error, "error returned when handling message. err: "+err.Error())
		}
		return err
	})
}
