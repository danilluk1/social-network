package gapi

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func grpcLogger(
	logger *zap.Logger,
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logFields := []zap.Field{
		zap.String("protocol", "grpc"),
		zap.String("method", info.FullMethod),
		zap.Int("status_code", int(statusCode)),
		zap.String("status_text", statusCode.String()),
		zap.Duration("duration", duration),
	}

	if err != nil {
		logFields = append(logFields, zap.Error(err))
	}

	logger.Info("received a gRPC Request", logFields...)
	return result, err
}

func GrpcLoggerWrapper(logger *zap.Logger) func(context.Context, interface{}, *grpc.UnaryServerInfo, grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		return grpcLogger(logger, ctx, req, info, handler)
	}
}
