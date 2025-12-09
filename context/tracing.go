package context

import (
	"net/http"

	"golang.org/x/net/context"

	"swallow-supplier/utils"
	utilsContext "swallow-supplier/utils/context"
)

const (
	traceIDLen = 15

	// TraceIDHeader header for Trace ID
	TraceIDHeader = "X-Trace-ID"

	// CtxLabelTraceID context label for trace id
	CtxLabelTraceID = "trace_id"
)

// TraceIDExtractor is a go-kit before handler that extracts a trace ID from HTTP headers, or creates a new one.
func TraceIDExtractor(ctx context.Context, r *http.Request) context.Context {
	traceID := r.Header.Get(TraceIDHeader)
	if traceID != "" {
		return ctxWithTraceID(ctx, traceID)
	}
	return ctxWithTraceID(ctx, utils.GenerateUUID("", true))
}

// TraceIDSetter is a go-kit after handler that sets the context trace ID into a HTTP header.
func TraceIDSetter(ctx context.Context, w http.ResponseWriter) context.Context {
	traceID := CtxTraceID(ctx)
	w.Header().Set(TraceIDHeader, traceID)
	return ctxWithTraceID(ctx, traceID)
}

func ctxWithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, CtxLabelTraceID, traceID)
}

// CtxTraceID returns traceId from the context
func CtxTraceID(ctx context.Context) string {
	return utilsContext.EnsureString(ctx, CtxLabelTraceID)
}
