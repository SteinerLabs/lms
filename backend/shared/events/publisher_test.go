package events

import (
	"context"
	"testing"

	"github.com/nats-io/nats.go"
)

// mockJetStream implements a mock NATS JetStream for testing
type mockJetStream struct {
	nats.JetStreamContext
	messages chan *nats.Msg
}

func (m *mockJetStream) Publish(subject string, data []byte, opts ...nats.PubOpt) (*nats.PubAck, error) {
	m.messages <- &nats.Msg{
		Subject: subject,
		Data:    data,
	}
	return &nats.PubAck{}, nil
}

func (m *mockJetStream) Subscribe(subject string, cb nats.MsgHandler, opts ...nats.SubOpt) (*nats.Subscription, error) {
	go func() {
		for msg := range m.messages {
			cb(msg)
		}
	}()
	return &nats.Subscription{}, nil
}

func TestPublisher_Publish(t *testing.T) {
	js := &mockJetStream{
		messages: make(chan *nats.Msg, 1),
	}

	publisher := NewPublisher(js, "test-service")

	type testData struct {
		Message string `json:"message"`
	}

	event := Event[any]{
		Type: "test.event",
		Data: testData{Message: "test"},
	}

	ctx := context.Background()
	err := publisher.Publish(ctx, event)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(js.messages) != 1 {
		t.Errorf("Expected 1 published message, got %d", len(js.messages))
	}

	if message := <-js.messages; message.Subject != "test.event" {
		t.Error("Expected message to be published to 'test.event'")
	}
}
