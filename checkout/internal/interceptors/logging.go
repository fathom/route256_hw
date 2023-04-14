package interceptors

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		logger.Info(
			fmt.Sprintf("gRPC request: %v", req),
			zap.String("full_method", info.FullMethod),
		)

		res, err := handler(ctx, req)
		if err != nil {
			logger.Warn(
				fmt.Sprintf("gRPC error: %v", err),
				zap.String("full_method", info.FullMethod),
			)
			return nil, err
		}

		logger.Info(
			fmt.Sprintf("gRPC response: %v", res),
			zap.String("full_method", info.FullMethod),
		)
		return res, nil
	}
}
