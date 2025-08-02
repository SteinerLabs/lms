package event

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/SteinerLabs/lms/backend/services/auth/internal/config"
	"github.com/SteinerLabs/lms/backend/shared/events"
)

// Publisher defines the interface for publishing events
type Publisher interface {
	// Publish publishes an event
	Publish(ctx context.Context, event *events.Event) error

	// Close closes the publisher
	Close() error
}

// KafkaPublisher implements the Publisher interface for Kafka
type KafkaPublisher struct {
	config *config.Config
	// In a real implementation, this would include a Kafka producer
	// producer *kafka.Producer
}

// NewKafkaPublisher creates a new Kafka publisher
func NewKafkaPublisher(cfg *config.Config) (*KafkaPublisher, error) {
	// In a real implementation, this would create a Kafka producer
	// producer, err := kafka.NewProducer(&kafka.ConfigMap{
	//     "bootstrap.servers": strings.Join(cfg.Kafka.Brokers, ","),
	// })
	// if err != nil {
	//     return nil, fmt.Errorf("failed to create Kafka producer: %w", err)
	// }

	return &KafkaPublisher{
		config: cfg,
		// producer: producer,
	}, nil
}

// Publish publishes an event to Kafka
func (p *KafkaPublisher) Publish(ctx context.Context, event *events.Event) error {
	// Marshal the event to JSON
	eventJSON, err := event.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	// In a real implementation, this would publish the event to Kafka
	// err = p.producer.Produce(&kafka.Message{
	//     TopicPartition: kafka.TopicPartition{Topic: &p.config.Kafka.Topic, Partition: kafka.PartitionAny},
	//     Value:          eventJSON,
	//     Key:            []byte(event.ID),
	// }, nil)
	// if err != nil {
	//     return fmt.Errorf("failed to produce message: %w", err)
	// }

	// For now, just log the event
	fmt.Printf("Publishing event: %s\n", string(eventJSON))

	return nil
}

// Close closes the Kafka publisher
func (p *KafkaPublisher) Close() error {
	// In a real implementation, this would close the Kafka producer
	// p.producer.Close()
	return nil
}

// MockPublisher implements the Publisher interface for testing
type MockPublisher struct {
	Events []*events.Event
}

// NewMockPublisher creates a new mock publisher
func NewMockPublisher() *MockPublisher {
	return &MockPublisher{
		Events: make([]*events.Event, 0),
	}
}

// Publish publishes an event to the mock publisher
func (p *MockPublisher) Publish(ctx context.Context, event *events.Event) error {
	p.Events = append(p.Events, event)
	return nil
}

// Close closes the mock publisher
func (p *MockPublisher) Close() error {
	return nil
}

// GetEvents returns the events published to the mock publisher
func (p *MockPublisher) GetEvents() []*events.Event {
	return p.Events
}

// GetEventsByType returns the events of a specific type published to the mock publisher
func (p *MockPublisher) GetEventsByType(eventType string) []*events.Event {
	events := make([]*events.Event, 0)
	for _, event := range p.Events {
		if event.Type == eventType {
			events = append(events, event)
		}
	}
	return events
}

// ClearEvents clears the events published to the mock publisher
func (p *MockPublisher) ClearEvents() {
	p.Events = make([]*events.Event, 0)
}

// PrintEvents prints the events published to the mock publisher
func (p *MockPublisher) PrintEvents() {
	for i, event := range p.Events {
		eventJSON, _ := json.MarshalIndent(event, "", "  ")
		fmt.Printf("Event %d: %s\n", i, string(eventJSON))
	}
}