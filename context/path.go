package context

import (
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"
)

// LabelKey represents the string label type
type LabelKey string

const (
	// CtxLabelRequestPath context request path label
	CtxLabelRequestPath LabelKey = "request_path"

	// CtxLabelRequestPathTemplate context request path template label
	CtxLabelRequestPathTemplate LabelKey = "request_path_template"

	// CtxLabelRequestURL context request url label
	CtxLabelRequestURL LabelKey = "request_url"
)

// RequestPathExtractor is a go-kit before interceptor that puts the request path in the request context.
func RequestPathExtractor(ctx context.Context, r *http.Request) context.Context {
	return CtxWithRequestPath(ctx, CtxLabelRequestPath, r.URL.EscapedPath())
}

// RequestPathTemplateExtractor is a go-kit before interceptor that puts the request path in the request context.
func RequestPathTemplateExtractor(ctx context.Context, r *http.Request) context.Context {
	router := mux.CurrentRoute(r)
	route, _ := router.GetPathTemplate()
	return CtxWithRequestPath(ctx, CtxLabelRequestPathTemplate, route)
}

// RequestURLExtractor is a go-kit before interceptor that puts the request url in the request context.
func RequestURLExtractor(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, CtxLabelRequestURL, r.URL)
}

// CtxWithRequestPath sets the path to context
func CtxWithRequestPath(ctx context.Context, label LabelKey, path string) context.Context {
	return context.WithValue(ctx, label, path)
}

// CtxRequestPath retrieves the path value from context
func CtxRequestPath(ctx context.Context, label LabelKey) string {
	v := ctx.Value(label)
	if v == nil {
		return ""
	}
	return v.(string)
}
