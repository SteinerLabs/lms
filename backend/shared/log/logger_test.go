package log

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		opts     []Option
		wantJSON bool
	}{
		{
			name:     "default logger",
			opts:     nil,
			wantJSON: false,
		},
		{
			name:     "JSON logger",
			opts:     []Option{WithJson()},
			wantJSON: true,
		},
		{
			name:     "text logger with custom prefix",
			opts:     []Option{WithPrefix("test"), WithoutJson()},
			wantJSON: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			opts := append(tt.opts, WithOutput(buf))
			logger := New(opts...)

			logger.Info("test message")
			if tt.wantJSON {
				if !json.Valid(buf.Bytes()) {
					t.Error("expected JSON output, got invalid JSON")
				}
			} else {
				if json.Valid(buf.Bytes()) {
					t.Error("expected text output, got JSON")
				}
			}
		})
	}
}

func TestLogLevels(t *testing.T) {
	tests := []struct {
		name      string
		logFn     func(*Logger, string, ...any)
		wantLevel string
	}{
		{"info", func(logger *Logger, s string, a ...any) {
			logger.Info(s, a...)
		}, "INFO"},
		{"warn", func(logger *Logger, s string, a ...any) {
			logger.Warn(s, a...)
		}, "WARN"},
		{"error", func(logger *Logger, s string, a ...any) {
			logger.Error(s, a...)
		}, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			logger := New(WithOutput(buf))

			tt.logFn(logger, "test message")
			output := buf.String()

			if !strings.Contains(output, tt.wantLevel) {
				t.Errorf("expected log level %s, got %s", tt.wantLevel, output)
			}
		})
	}
}

func TestTimeFormat(t *testing.T) {
	buf := &bytes.Buffer{}
	customFormat := "2006-01-02T15:04:05"
	logger := New(WithOutput(buf), WithTimeFormat(customFormat))

	logger.Info("test message")
	output := buf.String()

	// Parse the timestamp from the log output
	timestamp := strings.Split(output, " ")[0]
	_, err := time.Parse(customFormat, timestamp)
	if err != nil {
		t.Errorf("expected timestamp in format %s, got invalid format: %s", customFormat, timestamp)
	}
}

func TestWithSource(t *testing.T) {
	buf := &bytes.Buffer{}
	logger := New(WithOutput(buf), WithSource())

	logger.Info("test message")
	output := buf.String()

	if !strings.Contains(output, ".go:") {
		t.Error("expected source information in log output")
	}
}

func TestLoggerEnabled(t *testing.T) {
	tests := []struct {
		name        string
		level       slog.Level
		loggerLevel slog.Level
		want        bool
	}{
		{"debug enabled for debug logger", LevelDebug, LevelDebug, true},
		{"info disabled for error logger", LevelInfo, LevelError, false},
		{"error enabled for info logger", LevelError, LevelInfo, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(WithLevel(tt.loggerLevel))
			if got := logger.Enabled(context.Background(), tt.level); got != tt.want {
				t.Errorf("Logger.Enabled() = %v, want %v", got, tt.want)
			}
		})
	}
}
