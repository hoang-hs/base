package base

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.uber.org/fx"
)

func InitTracer(lc fx.Lifecycle) {
	cf := Get()
	if !cf.Tracer.Enabled {
		return
	}
	res, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			"",
			semconv.ServiceName(cf.Name),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", cf.Mode),
		),
	)
	opts := make([]sdktrace.TracerProviderOption, 0)
	opts = append(opts, sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()))

	if cf.Tracer.Jaeger.Active {
		expJaeger, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(cf.Tracer.Jaeger.Endpoint)))
		if err != nil {
			panic(err)
		}
		opts = append(opts, sdktrace.WithBatcher(expJaeger))
	}
	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}))
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return tp.Shutdown(ctx)
		},
	})
	return
}

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	traceId := ""
	if ctx.Value(TraceIdName) != nil {
		traceId = ctx.Value(TraceIdName).(string)
	}
	return traceId
}
