package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go/jetstream"
	"log"
	"time"
)

type HandlerFunc func(ctx context.Context, event Event[any]) error

type Consumer struct {
	js      jetstream.JetStream
	db      *sql.DB
	stream  string
	subject string
	durable string
}

func NewConsumer(js jetstream.JetStream, db *sql.DB, stream string, subject string, durable string) (*Consumer, error) {
	_, err := js.CreateStream(context.Background(), jetstream.StreamConfig{
		Name:     stream,
		Subjects: []string{subject},
		Storage:  jetstream.MemoryStorage,
		MaxAge:   24 * time.Hour, // Adjust retention period as needed
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create stream: %w", err)
	}

	return &Consumer{
		js:      js,
		db:      db,
		stream:  stream,
		subject: subject,
		durable: durable,
	}, nil
}

func (c *Consumer) Start(handler HandlerFunc) error {
	consumer, err := c.js.CreateConsumer(context.Background(), c.stream, jetstream.ConsumerConfig{
		Durable:       c.durable,
		Name:          c.durable,
		FilterSubject: c.subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return err
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		log.Printf("Received message: %s", msg.Data())
		var e Event[any]
		if err := json.Unmarshal(msg.Data(), &e); err != nil {
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

		err = c.markProcessed(e.ID)
		if err != nil {
			log.Printf("Failed to mark event as processed: %v", err)
		}
		err := msg.Ack()
		if err != nil {
			log.Printf("Failed to ack message: %v", err)
		}
	})
	if err != nil {
		return err
	}
	return nil
}

func (c *Consumer) hasProcessed(id string) (bool, error) {
	var exists bool
	err := c.db.QueryRow("SELECT EXISTS(SELECT 1 FROM processed_events WHERE event_id = $1)", id).Scan(&exists)
	return exists, err
}

func (c *Consumer) markProcessed(id string) error {
	log.Printf("Marking event as processed: %s", id)
	_, err := c.db.Exec("INSERT INTO processed_events(event_id) VALUES($1) ON CONFLICT DO NOTHING", id)
	return err
}
