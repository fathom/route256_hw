package interceptors

import (
	"context"
	"route256/loms/internal/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func Metrics(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	metrics.RequestCounter.Inc()

	timeStart := time.Now()

	res, err := handler(ctx, req)

	timeEnd := time.Since(timeStart)

	statusCode := status.Code(err).String()
	metrics.HistogramResponseTime.WithLabelValues(statusCode, info.FullMethod).Observe(timeEnd.Seconds())
	metrics.ResponseCounter.WithLabelValues(statusCode, info.FullMethod).Inc()

	return res, err
}
