package context

import (
	"context"
	"net/http"
	"swallow-supplier/utils"
)

// HeaderLabelKey represents the string label type
type HeaderLabelKey string

const (
	// RequestIDHeader represents X-Request-ID
	RequestIDHeader HeaderLabelKey = "X-Request-ID"
	// ChannelCodeHeader represents X-Tenant-Code
	ChannelCodeHeader HeaderLabelKey = "X-Channel-Code"
	// ApiKeyHeader represents X-Auth-Token
	ApiKeyHeader HeaderLabelKey = "X-Api-Key"
	// ForwardedForHeader represents X-Forwarded-For
	ForwardedForHeader HeaderLabelKey = "X-Forwarded-For"
	//Authorization represents AUB-Authorization
	Authorization HeaderLabelKey = "Authorization"
	// CtxLabelRequestID represents request_id label
	CtxLabelRequestID HeaderLabelKey = "Request_id"
	// CtxLabelChannelCode represents tenant_code label
	CtxLabelChannelCode HeaderLabelKey = "Channel_code"
	// CtxLabelApiKey represents auth_token label
	CtxLabelApiKey HeaderLabelKey = "Api_key"
	// CtxLabelForwardedFor represents forwarded_for label
	CtxLabelForwardedFor HeaderLabelKey = "forwarded_for"
	// CtxRemoteAddr represents remote_addr label
	CtxRemoteAddr HeaderLabelKey = "remote_addr"
	// CtxAubAuthorization represents Authorization label
	CtxAuthorization HeaderLabelKey = "Authorization"
)

// RequestIDHeaderExtractor is a go-kit before handler that extracts X-Request-ID.
func RequestIDHeaderExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxLabelRequestID, utils.GenerateUUID("GGT", true))
}

// ChannelCodeHeaderExtractor is a go-kit before handler that extracts X-Auth-Token.
func ChannelCodeHeaderExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxLabelChannelCode, r.Header.Get(string(ChannelCodeHeader)))
}

// ApiKeyHeaderExtractor is a go-kit before handler that extracts X-Channel-Code.
func ApiKeyHeaderExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxLabelApiKey, r.Header.Get(string(ApiKeyHeader)))
}

// ForwardedForHeaderExtractor is a go-kit before handler that extracts X-Forwarded-For.
func ForwardedForHeaderExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxLabelForwardedFor, r.Header.Get(string(ForwardedForHeader)))
}

// RemoteAddrExtractor is a go-kit before handler that extracts RemoteAddr.
func RemoteAddrExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxRemoteAddr, r.RemoteAddr)
}

func ctxWithHeader(ctx context.Context, label HeaderLabelKey, value string) context.Context {
	return context.WithValue(ctx, label, value)
}

// GGTAuthorizationExtractor is a go-kit before handler that extracts GGT-Authorization.
func GGTAuthorizationExtractor(ctx context.Context, r *http.Request) context.Context {
	return ctxWithHeader(ctx, CtxAuthorization, r.Header.Get(string(Authorization)))
}

// GetCtxHeader returns the header from context
func GetCtxHeader(ctx context.Context, label HeaderLabelKey) string {
	v := ctx.Value(label)
	if v == nil {
		return ""
	}
	return v.(string)
}
