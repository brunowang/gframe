package gflog

import (
	"context"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"time"
)

func UnaryEntryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	nowt := time.Now()
	defer func() {
		latency := time.Since(nowt)
		zaps := []zap.Field{
			zap.String("uri", info.FullMethod),
			zap.Duration("latency", latency),
		}
		Info(ctx, "grpc-entry", zaps...)
	}()

	return handler(ctx, req)
}

func StreamEntryInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	nowt := time.Now()
	defer func() {
		latency := time.Since(nowt)
		zaps := []zap.Field{
			zap.String("uri", info.FullMethod),
			zap.Duration("latency", latency),
		}
		Info(ctx, "grpc-entry", zaps...)
	}()

	return handler(srv, ss)
}
