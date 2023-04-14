package interceptors

import (
	"context"
	"route256/checkout/internal/metrics"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func ClientInterceptor(ctx context.Context, method string, req interface{}, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	timeStart := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)

	timeEnd := time.Since(timeStart)
	statusCode := status.Code(err).String()
	metrics.HistogramClientTime.WithLabelValues(statusCode, method).Observe(timeEnd.Seconds())

	return err
}
