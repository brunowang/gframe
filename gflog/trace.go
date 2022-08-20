package gflog

import "context"

type TraceKey struct{}

func GetTraceId(ctx context.Context) string {
	traceId, ok := ctx.Value(TraceKey{}).(string)
	if !ok || traceId == "" {
		return ""
	}
	return traceId
}
