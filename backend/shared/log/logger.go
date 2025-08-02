package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
)

const (
	LevelDebug = slog.LevelDebug
	LevelInfo  = slog.LevelInfo
	LevelWarn  = slog.LevelWarn
	LevelError = slog.LevelError
	LevelFatal = slog.Level(12)
)

type contextKey string

const (
	defaultLogLevel = slog.LevelInfo
	logWithFields   = contextKey("logWithFields")
)

// Logger represents a custom logger implementation utilizing slog package.
type Logger struct {
	*slog.Logger
	closer io.Closer // for cleanup of underlying resources
}

var (
	levelNames = map[slog.Leveler]string{
		LevelFatal: "FATAL",
	}
)

// Close properly closes the logger and its resources
func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// New creates a new Logger with all provided Option functions applied
func New(opts ...Option) *Logger {
	opt := *defaultConfig

	for _, configFn := range opts {
		configFn(&opt)
	}

	if opt.Json {
		h := slog.NewJSONHandler(opt.Output, &slog.HandlerOptions{
			AddSource: opt.AddSource,
			Level:     opt.Level,
			ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
				if a.Key == "source" {
					src := a.Value.Any().(*slog.Source)
					return slog.String("source", fmt.Sprintf("%v:%v", src.File, src.Line))
				}
				if a.Key == "level" {
					if lvl, ok := a.Value.Any().(slog.Level); ok {
						if name, exists := levelNames[lvl]; exists {
							return slog.String("level", name)
						}
					}
				}

				return a
			},
		})
		l := slog.New(h)
		if opt.Prefix != "" {
			l = l.With("prefix", opt.Prefix)
		}
		slog.SetDefault(l)
		return &Logger{Logger: l}
	} else {
		h := newHandler(&opt)
		l := slog.New(h)
		slog.SetDefault(l)
		return &Logger{Logger: l}
	}
}

// Fatal logs an error message with the specified msg and the provided args as slog.Attr using the receiver
// After the message is logged, the program exits with os.Exit(1)
func (l *Logger) Fatal(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, msg, args...)
	os.Exit(1)
}

// WithFields returns a new Logger instance with additional fields added to it.
// The additional fields are specified in the `fields` parameter as a map[string]any.
// The method iterates over the `fields` map and calls the `With` method on the original Logger,
// passing each key-value pair from the `fields` map. Finally, a new Logger instance is created
// using the updated Logger returned from the last `With` call, and this new Logger is returned.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	logger := l.Logger
	for k, v := range fields {
		logger = l.With(k, v)
	}
	return &Logger{
		logger,
		l,
	}
}

// WithField returns a new Logger instance with an additional field added to it.
// The additional field is specified by the `key` and `value` parameters.
// It calls the `With` method on the original Logger, passing the `key` and `value`.
// Finally, a new Logger instance is created using the updated Logger returned from the `With` call, and this new Logger is returned.
func (l *Logger) WithField(key string, value any) *Logger {
	return &Logger{
		l.With(key, value),
		l,
	}
}

// WithContext returns a new Logger instance with additional fields added to it based on the provided context.
// It retrieves the fields from the context using the logWithFields key and iterates over them, calling the `With` method on the original Logger,
// passing each key-value pair from the retrieved fields. Finally, it returns a new Logger instance using the updated Logger
// returned from the last `With` call.
func (l *Logger) WithContext(ctx context.Context) *Logger {
	ctxValue := ctx.Value(logWithFields)
	if ctxValue == nil {
		return l
	}
	fields := ctxValue.(map[string]any)
	logger := l.Logger
	contextFields := []any{}
	for k, v := range fields {
		fmt.Println(k, v)
		contextFields = append(contextFields, k, v)
	}
	logger = l.With(contextFields...)
	return &Logger{
		logger,
		l,
	}
}
