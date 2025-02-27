package observe

import (
	"context"
	"github.com/hoang-hs/base/src/common/log"
	"github.com/hoang-hs/base/src/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"time"
)

func InitTrace(lc fx.Lifecycle, cf *configs.Config, ctx context.Context) {
	if !cf.Observe.Trace.Enabled {
		return
	}
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(
		otlptracehttp.WithEndpoint(cf.Observe.Trace.OtlpExporter.Endpoint),
		otlptracehttp.WithInsecure(),
	))
	if err != nil {
		log.Fatal("Failed to create trace exporter", zap.Error(err))
	}
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceName(cf.Server.Name),
			semconv.ServiceVersion("v0.1.0"),
		),
	)

	traceProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(3*time.Second),
		),
		trace.WithResource(res),
		trace.WithSampler(trace.TraceIDRatioBased(cf.Observe.Trace.SampleRate)),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return traceProvider.Shutdown(ctx)
		},
	})
	return
}
