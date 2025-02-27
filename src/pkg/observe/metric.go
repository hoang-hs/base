package observe

import (
	"context"
	"github.com/Depado/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/hoang-hs/base/src/configs"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.uber.org/fx"
)

func NewMetric(lc fx.Lifecycle, cfg *configs.Config, r *gin.Engine) {
	if !cfg.Observe.Metric.Enabled {
		return
	}
	otelPromExporter, err := prometheus.New()
	if err != nil {
		panic(err)
	}
	meterProvider := metric.NewMeterProvider(metric.WithReader(otelPromExporter))
	otel.SetMeterProvider(meterProvider)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return meterProvider.Shutdown(ctx)
		},
	})
	// collect all the metrics
	p := ginprom.New(
		ginprom.Engine(r),
		ginprom.Subsystem(cfg.Server.Name),
		ginprom.Path("/metrics"),
	)
	r.Use(p.Instrument())
}
