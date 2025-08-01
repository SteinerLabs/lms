package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// Event represents a generic event in the system
type Event struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Source     string                 `json:"source"`
	Time       time.Time              `json:"time"`
	Data       interface{}            `json:"data"`
	Metadata   EventMetadata          `json:"metadata"`
}

// EventMetadata contains metadata about the event
type EventMetadata struct {
	CorrelationID string `json:"correlation_id"`
	CausationID   string `json:"causation_id"`
}

// NewEvent creates a new event with the given type, source, and data
func NewEvent(eventType, source string, data interface{}, correlationID, causationID string) *Event {
	if correlationID == "" {
		correlationID = uuid.New().String()
	}
	if causationID == "" {
		causationID = correlationID
	}

	return &Event{
		ID:     uuid.New().String(),
		Type:   eventType,
		Source: source,
		Time:   time.Now().UTC(),
		Data:   data,
		Metadata: EventMetadata{
			CorrelationID: correlationID,
			CausationID:   causationID,
		},
	}
}

// Marshal marshals the event to JSON
func (e *Event) Marshal() ([]byte, error) {
	return json.Marshal(e)
}

// Unmarshal unmarshals the event from JSON
func Unmarshal(data []byte, dataType interface{}) (*Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}

	// Unmarshal the data field into the provided type
	dataBytes, err := json.Marshal(event.Data)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(dataBytes, dataType); err != nil {
		return nil, err
	}

	event.Data = dataType
	return &event, nil
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
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	LogoutAt  time.Time `json:"logout_at"`
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