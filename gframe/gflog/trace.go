package gflog

import (
	"context"
	"github.com/brunowang/gframe/gfid"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

type TraceKey struct{}

func GetTraceID(ctx context.Context) string {
	traceID, ok := ctx.Value(TraceKey{}).(string)
	if !ok || traceID == "" {
		return ""
	}
	return traceID
}

func TraceCtx(ctx context.Context, traceIDs ...string) context.Context {
	if len(traceIDs) > 0 {
		return context.WithValue(ctx, TraceKey{}, strings.Join(traceIDs, "=>"))
	}
	return context.WithValue(ctx, TraceKey{}, gfid.GenID())
}

func UnaryTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var traceID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if trace, ok := md["trace_id"]; ok && len(trace) > 0 && trace[0] != "" {
			traceID = trace[0]
		}
		if traceID == "" {
			if trace, ok := md["uuid"]; ok && len(trace) > 0 && trace[0] != "" {
				traceID = trace[0]
			}
		}
	}
	if traceID == "" {
		traceID = gfid.GenID()
	}
	ctx = context.WithValue(ctx, TraceKey{}, traceID)
	return handler(ctx, req)
}

func StreamTraceInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := ss.Context()
	var traceID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if trace, ok := md["trace_id"]; ok && len(trace) > 0 && trace[0] != "" {
			traceID = trace[0]
		}
		if traceID == "" {
			if trace, ok := md["uuid"]; ok && len(trace) > 0 && trace[0] != "" {
				traceID = trace[0]
			}
		}
	}
	if traceID == "" {
		traceID = gfid.GenID()
	}
	ctx = context.WithValue(ctx, TraceKey{}, traceID)
	ssWrapper := grpc_middleware.WrapServerStream(ss)
	ssWrapper.WrappedContext = ctx
	return handler(srv, ssWrapper)
}
