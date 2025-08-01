package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/nats-io/nats.go"
	"log"
)

type HandlerFunc func(ctx context.Context, event Event[any]) error

type Consumer struct {
	js      nats.JetStreamContext
	db      *sql.DB
	subject string
	durable string
}

func NewConsumer(js nats.JetStreamContext, db *sql.DB, subject string, durable string) *Consumer {
	return &Consumer{
		js:      js,
		db:      db,
		subject: subject,
		durable: durable,
	}
}

func (c *Consumer) Start(handler HandlerFunc) error {
	_, err := c.js.Subscribe(c.subject, func(msg *nats.Msg) {
		var e Event[any]
		if err := json.Unmarshal(msg.Data, &e); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			err := msg.Nak()
			if err != nil {
				log.Printf("Failed to nak message: %v", err)
			}
			return
		}

		processed, _ := c.hasProcessed(e.ID)
		if processed {
			err := msg.Ack()
			if err != nil {
				log.Printf("Failed to ack message: %v", err)
			}
			return
		}

		ctx := WithEventContext(context.Background(), e)

		if err := handler(ctx, e); err != nil {
			log.Printf("Failed to handle event: %v", err)
			err := msg.Nak()
			if err != nil {
				log.Printf("Failed to nak message: %v", err)
			}
			return
		}

		_ = c.markProcessed(e.ID)
		err := msg.Ack()
		if err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	}, nats.Durable(c.durable), nats.ManualAck())
	return err
}

func (c *Consumer) hasProcessed(id string) (bool, error) {
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_events WHERE event_id = $1)", id).Scan(&exists)
	return exists, err
}

func (c *Consumer) markProcessed(id string) error {
	_, err := c.db.Exec("INSERT INTO processed_events(event_id) VALUES($1) ON CONFLICT DO NOTHING", id)
	return err
}
