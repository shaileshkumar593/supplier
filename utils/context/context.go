package context

import (
	"context"
)

type (
	consumerKey string
	timezoneKey string
)

const (
	// CtxLabelTimezone context timezone label
	CtxLabelTimezone consumerKey = "timezone"

	// CtxLabelConsumer context consumer key label
	CtxLabelConsumer timezoneKey = "consumer"
)

// EnsureString extracts a string from the context, returns "" if not found.
func EnsureString(ctx context.Context, label string) string {
	v := ctx.Value(label)
	if v == nil {
		return ""
	}
	return v.(string)
}

// EnsureInt64 extracts an int64 from the context, returns 0 if not found.
func EnsureInt64(ctx context.Context, label string) int64 {
	v := ctx.Value(label)
	if v == nil {
		return 0
	}
	return v.(int64)
}
