package tracing

import (
	"route256/loms/internal/logger"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"go.uber.org/zap"
)

func Init(jaegerHost, serviceName string) {
	logger.Info(
		"init tracing",
		zap.String("service_name", serviceName),
		zap.String("jaeger_host", jaegerHost),
	)

	cfg := &config.Configuration{
		ServiceName: serviceName,
		Sampler: &config.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: jaegerHost,
		},
	}

	_, err := cfg.InitGlobalTracer(serviceName)

	if err != nil {
		logger.Fatal("Cannot init tracing", zap.Error(err))
	}
}
