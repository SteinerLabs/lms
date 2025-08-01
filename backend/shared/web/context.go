package web

import (
	"context"
	"fmt"
	"time"
)

type ctxKey int

const key ctxKey = 1

type Values struct {
	TraceID    string
	Now        time.Time
	StatusCode int
}

func GetValues(ctx context.Context) (*Values, error) {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return nil, fmt.Errorf("no values in context")
	}
	return v, nil
}

func SetStatusCode(ctx context.Context, statusCode int) error {
	v, ok := ctx.Value(key).(*Values)
	if !ok {
		return fmt.Errorf("value missing from context")
	}
	v.StatusCode = statusCode
	return nil
}
