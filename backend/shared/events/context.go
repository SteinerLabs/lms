package events

import "context"

type ctxKey int

const (
	ctxTraceIDKey ctxKey = iota
	ctxCorrelationIDKey
	ctxCausationIDKey
)

func WithEventContext(ctx context.Context, event Event[any]) context.Context {
	ctx = context.WithValue(ctx, ctxTraceIDKey, event.TraceID)
	ctx = context.WithValue(ctx, ctxCorrelationIDKey, event.CorrelationID)
	ctx = context.WithValue(ctx, ctxCausationIDKey, event.CausationID)
	return ctx
}

func TraceIDFromContext(ctx context.Context) string {
	if v := ctx.Value(ctxTraceIDKey); v != nil {
		return v.(string)
	}
	return ""
}
