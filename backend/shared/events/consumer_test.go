package events

import (
	"context"
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"
	"time"

	"github.com/nats-io/nats.go"
)

// Create a new mock DB
func newMockDB() *sql.DB {
	db, mock, err := sqlmock.New()
	if err != nil {
		panic(err)
	}

	// Set up expectations for the queries we'll use
	mock.ExpectQuery("SELECT EXISTS").WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
	mock.ExpectExec("INSERT INTO processed_events").WillReturnResult(sqlmock.NewResult(1, 1))

	return db
}

func TestConsumer_Start(t *testing.T) {
	// Create a mock NATS JetStream
	mockJS := &mockJetStream{
		messages: make(chan *nats.Msg, 1),
	}

	// Create a mock DB
	db := newMockDB()
	defer db.Close()

	consumer := NewConsumer(mockJS, db, "test.subject", "test-durable")

	// Create a test handler
	handlerCalled := false
	handler := func(ctx context.Context, event Event[any]) error {
		handlerCalled = true
		return nil
	}

	// Start the consumer in a goroutine
	errChan := make(chan error, 1)
	go func() {
		errChan <- consumer.Start(handler)
	}()

	// Create and send a test event
	testEvent := Event[any]{
		ID:         "test-id",
		Type:       "test.event",
		OccurredAt: time.Now(),
	}
	eventData, _ := json.Marshal(testEvent)
	msg := &nats.Msg{
		Subject: "test.subject",
		Data:    eventData,
	}

	// Send the message to the mock NATS
	mockJS.messages <- msg

	// Wait a bit for processing
	time.Sleep(100 * time.Millisecond)

	// Verify that the handler was called
	if !handlerCalled {
		t.Error("Expected handler to be called")
	}
}
