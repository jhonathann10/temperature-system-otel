package providerotel

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func InitProvider() *sdktrace.TracerProvider {
	endpoint := "http://localhost:9411/api/v2/spans"
	exporter, err := zipkin.New(endpoint)
	if err != nil {
		log.Fatalf("failed to create Zipkin exporter: %v", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("temperature-system-otel"),
			semconv.ServiceNameKey.String("temperature-system-otel"),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp
}
