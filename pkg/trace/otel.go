package trace

import (
	"context"
	"fmt"
	"io"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

type ExporterType int

const (
	Undefined ExporterType = iota
	NOP
	Stdout
	Datadog
)

const tracerName = "github.com/ca-risken/common/pkg/trace"

func (t ExporterType) String() string {
	switch t {
	case NOP:
		return "nop"
	case Stdout:
		return "stdout"
	case Datadog:
		return "datadog"
	default:
		return "undefined"
	}
}

func GetExporterType(typeString string) ExporterType {
	switch typeString {
	case "nop":
		return NOP
	case "stdout":
		return Stdout
	case "datadog":
		return Datadog
	default:
		return Undefined
	}
}

type Config struct {
	Namespace   string
	ServiceName string
	Environment string
	ExporterType
}

func (c *Config) GetFullServiceName() string {
	return fmt.Sprintf("%s.%s", c.Namespace, c.ServiceName)
}

func Init(ctx context.Context, config *Config) (*trace.TracerProvider, error) {
	var exporter trace.SpanExporter
	var err error
	switch config.ExporterType {
	case NOP:
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(io.Discard))
	case Stdout:
		exporter, err = stdouttrace.New()
	case Datadog:
		client := otlptracehttp.NewClient()
		exporter, err = otlptrace.New(ctx, client)
	default:
		// fallback to nop exporter
		exporter, err = stdouttrace.New(
			stdouttrace.WithWriter(io.Discard))
	}
	if err != nil {
		return nil, err
	}

	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(config.GetFullServiceName()),
			attribute.String("environment", config.Environment)))
	if err != nil {
		return nil, err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(r),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
