package telemetry

import (
	"context"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// Provider хранит провайдер трейсов для корректного shutdown.
type Provider struct {
	tp *sdktrace.TracerProvider
}

// InitOTLP инициализирует OTLP trace провайдер.
// endpoint — адрес коллектора (например "jaeger:4318" или "localhost:4318").
func InitOTLP(ctx context.Context, endpoint, serviceName string, log *slog.Logger) (*Provider, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	// ----- Traces -----
	traceExp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExp, sdktrace.WithBatchTimeout(5*time.Second)),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	log.Info("OTLP telemetry initialized",
		slog.String("endpoint", endpoint),
		slog.String("service", serviceName),
	)

	return &Provider{tp: tp}, nil
}

// Shutdown корректно завершает работу провайдера.
func (p *Provider) Shutdown(ctx context.Context) error {
	if p == nil {
		return nil
	}
	if p.tp != nil {
		return p.tp.Shutdown(ctx)
	}
	return nil
}
