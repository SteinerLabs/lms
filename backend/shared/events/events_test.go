// events_test.go
package events

import (
	"testing"
)

func TestNewEvent(t *testing.T) {
	type testData struct {
		Message string
	}

	tests := []struct {
		name          string
		eventType     string
		source        string
		data          testData
		correlationID string
		causationID   string
		traceID       string
	}{
		{
			name:      "with all IDs empty",
			eventType: "test.created",
			source:    "test-service",
			data:      testData{Message: "test"},
		},
		{
			name:          "with correlation ID",
			eventType:     "test.created",
			source:        "test-service",
			data:          testData{Message: "test"},
			correlationID: "corr123",
		},
		{
			name:          "with all IDs",
			eventType:     "test.created",
			source:        "test-service",
			data:          testData{Message: "test"},
			correlationID: "corr123",
			causationID:   "cause123",
			traceID:       "trace123",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := NewEvent(tt.eventType, tt.source, tt.data, tt.correlationID, tt.causationID, tt.traceID)

			if event.Type != tt.eventType {
				t.Errorf("Want event type %q, got %q", tt.eventType, event.Type)
			}

			if event.Source != tt.source {
				t.Errorf("Want source %q, got %q", tt.source, event.Source)
			}

			if event.ID == "" {
				t.Error("Event ID should not be empty")
			}

			if event.CorrelationID == "" {
				t.Error("CorrelationID should not be empty")
			}

			if event.CausationID == "" {
				t.Error("CausationID should not be empty")
			}

			if event.TraceID == "" {
				t.Error("TraceID should not be empty")
			}

			if event.Data != tt.data {
				t.Errorf("Want data %v, got %v", tt.data, event.Data)
			}

			if event.OccurredAt.IsZero() {
				t.Error("OccurredAt should not be zero")
			}
		})
	}
}
