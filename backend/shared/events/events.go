package events

import (
	"time"

	"github.com/google/uuid"
)

// Event represents a generic event in the system
type Event[T any] struct {
	ID            string    `json:"id"`
	TraceID       string    `json:"trace_id"`
	CorrelationID string    `json:"correlation_id"`
	CausationID   string    `json:"causation_id"`
	Source        string    `json:"source"` // Example: auth-service
	Type          string    `json:"type"`
	OccurredAt    time.Time `json:"occurred_at"`
	Data          T         `json:"data"`
}

// NewEvent creates a new event with the given type, source, and data
func NewEvent[T any](eventType string, source string, data T, correlationID, causationID string, traceId string) *Event[T] {
	if correlationID == "" {
		correlationID = uuid.New().String()
	}
	if causationID == "" {
		causationID = correlationID
	}
	if traceId == "" {
		traceId = correlationID
	}

	return &Event[T]{
		ID:            uuid.New().String(),
		Type:          eventType,
		OccurredAt:    time.Now().UTC(),
		Data:          data,
		CorrelationID: correlationID,
		CausationID:   causationID,
		TraceID:       traceId,
		Source:        source,
	}
}

// User Events

// UserCreatedEvent represents a user.created event
type UserCreatedEvent struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
}

// UserUpdatedEvent represents a user.updated event
type UserUpdatedEvent struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserDeletedEvent represents a user.deleted event
type UserDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UserLoggedInEvent represents a user.login event
type UserLoggedInEvent struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	IP        string    `json:"ip"`
	UserAgent string    `json:"user_agent"`
	LoginAt   time.Time `json:"login_at"`
}

// UserLoggedOutEvent represents a user.logout event
type UserLoggedOutEvent struct {
	ID       string    `json:"id"`
	Email    string    `json:"email"`
	LogoutAt time.Time `json:"logout_at"`
}

// Course Events

// CourseCreatedEvent represents a course.created event
type CourseCreatedEvent struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	InstructorIDs []string  `json:"instructor_ids"`
	CreatedAt     time.Time `json:"created_at"`
}

// CoursePublishedEvent represents a course.published event
type CoursePublishedEvent struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	InstructorIDs []string  `json:"instructor_ids"`
	PublishedAt   time.Time `json:"published_at"`
}

// CourseUpdatedEvent represents a course.updated event
type CourseUpdatedEvent struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CourseDeletedEvent represents a course.deleted event
type CourseDeletedEvent struct {
	ID        string    `json:"id"`
	DeletedAt time.Time `json:"deleted_at"`
}

// UserEnrolledEvent represents a course.enrollment.created event
type UserEnrolledEvent struct {
	ID         string    `json:"id"`
	CourseID   string    `json:"course_id"`
	UserID     string    `json:"user_id"`
	EnrolledAt time.Time `json:"enrolled_at"`
}

// EnrollmentCompletedEvent represents a course.enrollment.completed event
type EnrollmentCompletedEvent struct {
	ID          string    `json:"id"`
	CourseID    string    `json:"course_id"`
	UserID      string    `json:"user_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// Progress Events

// ProgressUpdatedEvent represents a progress.updated event
type ProgressUpdatedEvent struct {
	UserID          string    `json:"user_id"`
	CourseID        string    `json:"course_id"`
	ModuleID        string    `json:"module_id"`
	LessonID        string    `json:"lesson_id"`
	PercentComplete float64   `json:"percent_complete"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// LessonCompletedEvent represents a progress.lesson.completed event
type LessonCompletedEvent struct {
	UserID      string    `json:"user_id"`
	CourseID    string    `json:"course_id"`
	ModuleID    string    `json:"module_id"`
	LessonID    string    `json:"lesson_id"`
	CompletedAt time.Time `json:"completed_at"`
}

// QuizCompletedEvent represents a progress.quiz.completed event
type QuizCompletedEvent struct {
	UserID      string    `json:"user_id"`
	CourseID    string    `json:"course_id"`
	QuizID      string    `json:"quiz_id"`
	Score       float64   `json:"score"`
	Passed      bool      `json:"passed"`
	CompletedAt time.Time `json:"completed_at"`
}

// AssignmentSubmittedEvent represents a progress.assignment.submitted event
type AssignmentSubmittedEvent struct {
	UserID       string    `json:"user_id"`
	CourseID     string    `json:"course_id"`
	AssignmentID string    `json:"assignment_id"`
	SubmittedAt  time.Time `json:"submitted_at"`
}

// AchievementEarnedEvent represents a progress.achievement.earned event
type AchievementEarnedEvent struct {
	UserID          string    `json:"user_id"`
	AchievementID   string    `json:"achievement_id"`
	AchievementName string    `json:"achievement_name"`
	EarnedAt        time.Time `json:"earned_at"`
}

// Billing Events

// PaymentCompletedEvent represents a billing.payment.completed event
type PaymentCompletedEvent struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Amount      int       `json:"amount"`
	Currency    string    `json:"currency"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completed_at"`
}

// PaymentFailedEvent represents a billing.payment.failed event
type PaymentFailedEvent struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	Amount       int       `json:"amount"`
	Currency     string    `json:"currency"`
	ErrorMessage string    `json:"error_message"`
	FailedAt     time.Time `json:"failed_at"`
}

// Notification Events

// NotificationSentEvent represents a notification.sent event
type NotificationSentEvent struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Channels []string  `json:"channels"`
	SentAt   time.Time `json:"sent_at"`
}

// NotificationDeliveredEvent represents a notification.delivered event
type NotificationDeliveredEvent struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ChannelType string    `json:"channel_type"`
	DeliveredAt time.Time `json:"delivered_at"`
}

// NotificationReadEvent represents a notification.read event
type NotificationReadEvent struct {
	ID     string    `json:"id"`
	UserID string    `json:"user_id"`
	ReadAt time.Time `json:"read_at"`
}

// Example usage:
//
// // Create a new user created event
// userData := &UserCreatedEvent{
//     ID:        "123e4567-e89b-12d3-a456-426614174000",
//     Email:     "user@example.com",
//     FirstName: "John",
//     LastName:  "Doe",
//     CreatedAt: time.Now().UTC(),
// }
//
// // Create a new event with the user created data
// event := NewEvent("user.created", "auth-service", userData, "", "")
//
// // Marshal the event to JSON
// eventJSON, err := event.Marshal()
// if err != nil {
//     log.Fatalf("Failed to marshal event: %v", err)
// }
//
// // Publish the event to Kafka
// err = kafkaProducer.Publish("user-events", eventJSON)
// if err != nil {
//     log.Fatalf("Failed to publish event: %v", err)
// }
//
// // Consume the event from Kafka
// eventJSON := kafkaConsumer.Consume("user-events")
//
// // Unmarshal the event
// var userData UserCreatedEvent
// event, err := Unmarshal(eventJSON, &userData)
// if err != nil {
//     log.Fatalf("Failed to unmarshal event: %v", err)
// }
//
// // Use the event data
// fmt.Printf("User created: %s %s (%s)\n", userData.FirstName, userData.LastName, userData.Email)

/*
// event.go
package events

import (
	"encoding/json"
	"time"
)

// EventData is the interface that all event payloads must implement
type EventData interface {
	EventType() string
}

// Event represents the base event structure
type Event struct {
	ID            string          `json:"id"`
	Type          string          `json:"type"`
	TraceID       string          `json:"trace_id"`
	CorrelationID string          `json:"correlation_id"`
	CausationID   string          `json:"causation_id"`
	Source        string          `json:"source"`
	OccurredAt    time.Time       `json:"occurred_at"`
	Data          json.RawMessage `json:"data"`
}

// EventHandler handles a specific event type
type EventHandler interface {
	EventType() string
	Handle(ctx context.Context, data json.RawMessage) error
}

// TypedEventHandler is a type-safe wrapper for handling specific event types
type TypedEventHandler[T EventData] struct {
	handler func(ctx context.Context, data T) error
}

func NewTypedEventHandler[T EventData](h func(ctx context.Context, data T) error) *TypedEventHandler[T] {
	return &TypedEventHandler[T]{handler: h}
}

func (h *TypedEventHandler[T]) EventType() string {
	var t T
	return t.EventType()
}

func (h *TypedEventHandler[T]) Handle(ctx context.Context, data json.RawMessage) error {
	var typed T
	if err := json.Unmarshal(data, &typed); err != nil {
		return err
	}
	return h.handler(ctx, typed)
}

// HandlerRegistry manages event handlers
type HandlerRegistry struct {
	handlers map[string][]EventHandler
}

func NewHandlerRegistry() *HandlerRegistry {
	return &HandlerRegistry{
		handlers: make(map[string][]EventHandler),
	}
}

func (r *HandlerRegistry) Register(handler EventHandler) {
	eventType := handler.EventType()
	if r.handlers[eventType] == nil {
		r.handlers[eventType] = make([]EventHandler, 0)
	}
	r.handlers[eventType] = append(r.handlers[eventType], handler)
}

// publisher.go
type Publisher struct {
	js     nats.JetStreamContext
	source string
}

func (p *Publisher) Publish(ctx context.Context, data EventData) error {
	event := Event{
		ID:         uuid.New().String(),
		Type:       data.EventType(),
		Source:     p.source,
		OccurredAt: time.Now().UTC(),
	}

	// Set tracing info
	event.TraceID = TraceIDFromContext(ctx)
	if event.TraceID == "" {
		event.TraceID = uuid.New().String()
	}

	event.CorrelationID = event.TraceID
	if v := ctx.Value(ctxCorrelationIDKey); v != nil {
		event.CorrelationID = v.(string)
	}

	event.CausationID = event.CorrelationID
	if v := ctx.Value(ctxCausationIDKey); v != nil {
		event.CausationID = v.(string)
	}

	// Marshal the data
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	event.Data = dataBytes

	// Marshal the complete event
	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = p.js.Publish(event.Type, payload)
	return err
}

// consumer.go
type Consumer struct {
	js       nats.JetStreamContext
	db       *sql.DB
	registry *HandlerRegistry
	subject  string
	durable  string
}

func (c *Consumer) Start(ctx context.Context) error {
	_, err := c.js.Subscribe(c.subject, func(msg *nats.Msg) {
		var event Event
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Printf("Failed to unmarshal event: %v", err)
			_ = msg.Nak()
			return
		}

		if processed, _ := c.hasProcessed(event.ID); processed {
			_ = msg.Ack()
			return
		}

		handlers, exists := c.registry.handlers[event.Type]
		if !exists {
			log.Printf("No handlers registered for event type: %s", event.Type)
			_ = msg.Ack() // Ack anyway as we don't want to reprocess
			return
		}

		ctx = WithEventContext(ctx, event)

		for _, handler := range handlers {
			if err := handler.Handle(ctx, event.Data); err != nil {
				log.Printf("Handler failed: %v", err)
				_ = msg.Nak()
				return
			}
		}

		_ = c.markProcessed(event.ID)
		_ = msg.Ack()
	}, nats.Durable(c.durable), nats.ManualAck())

	return err
}

// Example usage:
type UserCreated struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func (e UserCreated) EventType() string {
	return "user.created"
}

func Example() {
	registry := NewHandlerRegistry()

	// Register a handler for UserCreated events
	registry.Register(NewTypedEventHandler(func(ctx context.Context, data UserCreated) error {
		// Handle the event
		log.Printf("User created: %s (%s)", data.UserID, data.Email)
		return nil
	}))

	publisher := NewPublisher(js, "user-service")
	consumer := NewConsumer(js, db, registry, "events.>", "user-consumer")

	// Publishing
	userData := UserCreated{
		UserID: "123",
		Email:  "user@example.com",
	}
	err := publisher.Publish(ctx, userData)

	// Start consuming
	consumer.Start(context.Background())
}
*/
