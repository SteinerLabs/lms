package events

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

type Publisher struct {
	js     nats.JetStreamContext
	source string // Example: auth-service
}

func NewPublisher(js nats.JetStreamContext, source string) *Publisher {
	return &Publisher{
		js:     js,
		source: source,
	}
}

func (p *Publisher) Publish(ctx context.Context, event Event[any]) error {
	traceID := TraceIDFromContext(ctx)
	if traceID != "" {
		traceID = uuid.New().String()
	}

	correlationID := traceID
	if v := ctx.Value(ctxCorrelationIDKey); v != nil {
		correlationID = v.(string)
	}

	causationID := correlationID
	if v := ctx.Value(ctxCausationIDKey); v != nil {
		causationID = v.(string)
	}

	event.TraceID = traceID
	event.CorrelationID = correlationID
	event.CausationID = causationID
	event.Source = p.source

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = p.js.Publish(event.Type, payload)
	return err
}
