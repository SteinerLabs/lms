package events

import (
	"context"
	"testing"
)

func TestEventContext(t *testing.T) {
	event := Event[any]{
		TraceID:       "trace123",
		CorrelationID: "corr123",
		CausationID:   "cause123",
	}

	ctx := context.Background()
	ctx = WithEventContext(ctx, event)

	tests := []struct {
		name     string
		getValue func(context.Context) string
		want     string
	}{
		{
			name:     "TraceID",
			getValue: TraceIDFromContext,
			want:     "trace123",
		},
		{
			name: "CorrelationID",
			getValue: func(ctx context.Context) string {
				v := ctx.Value(ctxCorrelationIDKey)
				if v == nil {
					return ""
				}
				return v.(string)
			},
			want: "corr123",
		},
		{
			name: "CausationID",
			getValue: func(ctx context.Context) string {
				v := ctx.Value(ctxCausationIDKey)
				if v == nil {
					return ""
				}
				return v.(string)
			},
			want: "cause123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.getValue(ctx); got != tt.want {
				t.Errorf("Want %q, got %q", tt.want, got)
			}
		})
	}
}
