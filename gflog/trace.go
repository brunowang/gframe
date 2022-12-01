package gflog

import (
	"context"
	"github.com/brunowang/gframe/gfid"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"strings"
)

const GinCtxKey = "__gin_ctx_key__"

type TraceKey struct{}

func GetTraceID(ctx context.Context) string {
	if ginCtx, ok := ctx.(*gin.Context); ok {
		ctx = GetCtxFromGin(ginCtx)
	}
	traceID, ok := ctx.Value(TraceKey{}).(string)
	if !ok || traceID == "" {
		return ""
	}
	return traceID
}

func GetCtxFromGin(c *gin.Context) context.Context {
	val, ok := c.Get(GinCtxKey)
	if !ok || val == nil {
		return context.TODO()
	}
	ctx, ok := val.(context.Context)
	if !ok || ctx == nil {
		return context.TODO()
	}
	return ctx
}

func TraceCtx(ctx context.Context, traceIDs ...string) context.Context {
	if len(traceIDs) > 0 {
		return context.WithValue(ctx, TraceKey{}, strings.Join(traceIDs, "=>"))
	}
	return context.WithValue(ctx, TraceKey{}, gfid.GenID())
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("trace_id")
		if traceID == "" {
			traceID = gfid.GenID()
		}

		ctx := context.TODO()
		if val, ok := c.Get(GinCtxKey); !ok || val == nil {
			ctx = context.WithValue(ctx, TraceKey{}, traceID)
		} else if c, ok := val.(context.Context); !ok {
			ctx = context.WithValue(ctx, TraceKey{}, traceID)
		} else {
			ctx = context.WithValue(c, TraceKey{}, traceID)
		}

		c.Set(GinCtxKey, ctx)
		c.Next()
	}
}

func UnaryTraceInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	var traceID string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if trace, ok := md["trace_id"]; ok && len(trace) > 0 && trace[0] != "" {
			traceID = trace[0]
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
	}
	if traceID == "" {
		traceID = gfid.GenID()
	}
	ctx = context.WithValue(ctx, TraceKey{}, traceID)
	ssWrapper := grpc_middleware.WrapServerStream(ss)
	ssWrapper.WrappedContext = ctx
	return handler(srv, ssWrapper)
}
