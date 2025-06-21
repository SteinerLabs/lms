package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strconv"
	"sync"
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

	defaultConfig = &Options{
		Output:     os.Stdout,
		AddSource:  true,
		Level:      defaultLogLevel,
		TimeFormat: "",
		Prefix:     "",
		Json:       false,
	}
)

// Close properly closes the logger and its resources
func (l *Logger) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// Option uses the functional options patter to configure a Logger
type Option func(*Options)

// WithLevel returns an Option that sets the log level of the Options struct to the provided level.
// The log level determines the verbosity of the log messages.
func WithLevel(level slog.Level) Option {
	return func(o *Options) {
		o.Level = level
	}
}

// WithSource returns an Option function that sets the AddSource field of the Options struct to true.
func WithSource() Option {
	return func(o *Options) {
		o.AddSource = true
	}
}

// WithoutSource returns an Option function that sets the AddSource field of the Options struct to false.
// This means that the source file and line number will not be included in log output.
func WithoutSource() Option {
	return func(o *Options) {
		o.AddSource = false
	}
}

// WithTimeFormat returns an Option that sets the time format of the Options struct
// to the provided format string. The time format is used when printing log messages
// that contain timestamps.
func WithTimeFormat(format string) Option {
	return func(o *Options) {
		o.TimeFormat = format
	}
}

// WithPrefix is an Option function that sets the prefix of the log messages
func WithPrefix(prefix string) Option {
	return func(o *Options) {
		o.Prefix = prefix
	}
}

// WithOutput returns an Option function that sets the Output to the provided io.Writer
func WithOutput(output io.Writer) Option {
	return func(o *Options) {
		o.Output = output
	}
}

// WithJson returns an Option function that enables Json-Logging
func WithJson() Option {
	return func(o *Options) {
		o.Json = true
	}
}

// WithoutJson returns an Option function that disables Json-Logging
func WithoutJson() Option {
	return func(o *Options) {
		o.Json = false
	}
}

// WithText returns an Option function that enables Text-Logging
var WithText = WithoutJson

type Options struct {
	Output     io.Writer
	AddSource  bool
	Level      slog.Leveler
	TimeFormat string
	Prefix     string
	Json       bool
}

type Handler struct {
	sync.RWMutex
	options    Options
	prefix     string
	preformat  string
	timeFormat string
	w          io.Writer
}

// newHandler creates a new instance of Handler with the provided options.
// It sets the output writer, options, time format, and prefix of the handler.
// Returns the new handler instance.
func newHandler(options *Options) *Handler {
	h := &Handler{w: options.Output}
	h.options = *options

	if options.TimeFormat != "" {
		h.timeFormat = options.TimeFormat
	} else {
		h.timeFormat = "2006/01/02 15:04:05"
	}

	h.prefix = options.Prefix
	return h
}

// Enabled checks if the Logger is enabled for a specific slog.Level
func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := defaultLogLevel
	if h.options.Level != nil {
		minLevel = h.options.Level.Level()
	}
	return level >= minLevel
}

// WithGroup creates a new Handler with a modified prefix by appending the provided name to the existing prefix.
func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{
		w:          h.w,
		options:    h.options,
		preformat:  h.preformat,
		timeFormat: h.timeFormat,
		prefix:     h.prefix + "." + name,
	}
}

// WithAttrs creates a new Handler with additional attributes appended to the preformat string.
func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	var buf []byte
	for _, a := range attrs {
		buf = h.appendAttr(buf, h.prefix, a)
	}
	return &Handler{
		w:          h.w,
		options:    h.options,
		prefix:     h.prefix,
		timeFormat: h.timeFormat,
		preformat:  h.preformat + string(buf),
	}
}

// Handle writes the Log Record to the specified output Writer in the required format.
// The Log Record consists of the timestamp, prefix, log level, source location (if enabled),
// message, preformat string, and any additional attributes.
// It acquires a lock on the Handler's mutex to ensure thread safety during writing.
// Returns any error that occurred during writing.
func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	buf := h.formatRecord(&r)

	h.Lock()
	defer h.Unlock()

	_, err := h.w.Write(buf)
	if err != nil {
		return fmt.Errorf("failed to write log record: %w", err)
	}

	if r.Level == LevelFatal {
		if f, ok := h.w.(interface{ Sync() error }); ok {
			if err := f.Sync(); err != nil {
				return fmt.Errorf("failed to sync log: %w", err)
			}
		}
	}

	return nil
}

// formatRecord formats a log record into a byte slice
func (h *Handler) formatRecord(r *slog.Record) []byte {
	h.RLock()
	defer h.RUnlock()

	var buf []byte

	// Time
	if !r.Time.IsZero() {
		buf = r.Time.AppendFormat(buf, h.timeFormat)
	}

	// Prefix
	if h.prefix != "" {
		buf = append(buf, ' ')
		buf = append(buf, h.prefix...)
	}

	// Level
	levelStr := r.Level.String()
	buf = append(buf, ' ')
	if _, ok := levelNames[r.Level]; ok {
		levelStr = levelNames[r.Level]
	}
	buf = append(buf, fmt.Sprintf("%-5s", levelStr)...)
	buf = append(buf, ' ')

	// Source
	if h.options.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		buf = append(buf, f.File...)
		buf = append(buf, ':')
		buf = strconv.AppendInt(buf, int64(f.Line), 10)
		buf = append(buf, ' ')
	}

	// Message
	buf = append(buf, ' ')
	buf = append(buf, r.Message...)
	buf = append(buf, h.preformat...)

	// Attributes
	r.Attrs(func(a slog.Attr) bool {
		buf = h.appendAttr(buf, h.prefix, a)
		return true
	})
	return append(buf, '\n')
}

// appendAttr appends the provided attribute to the buffer.
// If the attribute is empty, it returns the buffer as is.
// If the attribute value is not a group, it appends the key-value pair to the buffer.
// If the attribute value is a group, it recursively appends each attribute in the group to the buffer.
// The key is modified by adding the prefix to it, if specified.
// Returns the updated buffer.
func (h *Handler) appendAttr(buf []byte, prefix string, a slog.Attr) []byte {
	if a.Equal(slog.Attr{}) {
		return buf
	}
	if a.Value.Kind() != slog.KindGroup {
		buf = append(buf, ' ')
		buf = append(buf, a.Key...)
		buf = append(buf, '=')
		return fmt.Appendf(buf, "%v", a.Value.Any())
	}
	// Group
	if a.Key != "" {
		prefix += a.Key + "."
	}
	for _, a := range a.Value.Group() {
		buf = h.appendAttr(buf, prefix, a)
	}
	return buf
}

// New creates a new Logger with all provided Option functions applied
func New(opts ...Option) *Logger {
	opt := defaultConfig

	for _, configFn := range opts {
		configFn(opt)
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
		h := newHandler(opt)
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

// Fatalf logs an error message with the specified format and arguments using the receiver
// After the message is logged, the program exits with os.Exit(1)
func (l *Logger) Fatalf(msg string, args ...any) {
	l.Log(context.Background(), LevelFatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Infof logs an informational message with the specified format and arguments using the receiver
func (l *Logger) Infof(msg string, args ...any) {
	l.Info(fmt.Sprintf(msg, args...))
}

// Errorf logs an error message with the specified format and arguments using the receiver
func (l *Logger) Errorf(msg string, args ...any) {
	l.Error(fmt.Sprintf(msg, args...))
}

// Debugf logs a debug message with the specified format and arguments using the receiver
func (l *Logger) Debugf(msg string, args ...any) {
	l.Debug(fmt.Sprintf(msg, args...))
}

// Warnf logs a warning message with the specified format and arguments using the receiver
func (l *Logger) Warnf(msg string, args ...any) {
	l.Warn(fmt.Sprintf(msg, args...))
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
	for k, v := range fields {
		logger = l.With(k, v)
	}
	return &Logger{
		logger,
		l,
	}
}

// ContextWithValue adds a key-value pair to the context's logWithFields map.
// If the map is nil, it creates a new map and assigns it to the context.
// If the map already exists, it performs a shallow copy of the existing map,
// adds the key-value pair to the copy, and assigns it to the context.
// Lastly, it returns the updated context.
// This function is used to attach additional log fields to a specific context.
func ContextWithValue(ctx context.Context, key string, value any) context.Context {
	ctxValue := ctx.Value(logWithFields)
	if ctxValue == nil {
		return context.WithValue(ctx, logWithFields, map[string]any{key: value})
	}

	fields := shallowCopy(ctxValue.(map[string]interface{}))
	fields[key] = value

	return context.WithValue(ctx, logWithFields, fields)
}

// ContextWithValues creates a new context with additional key-value pairs specified in the fields parameter.
// If the supplied context already contains a value for the key "logWithFields", then the function merges the existing fields
// with the new fields using the "merge" function. Otherwise, a new logWithFields key-value pair is added to the context
// using the shallowCopy of the fields.
func ContextWithValues(ctx context.Context, fields map[string]any) context.Context {
	ctxValue := ctx.Value(logWithFields)

	if ctxValue == nil {
		return context.WithValue(ctx, logWithFields, shallowCopy(fields))
	}

	return context.WithValue(ctx, logWithFields, merge(ctxValue.(map[string]interface{}), fields))
}

// shallowCopy performs a shallow copy of the given map[string]any and
// returns a new map with the copied key-value pairs. If the given map is empty,
// it returns an empty map.
//
// Example:
//
//	m := map[string]any{"foo": 1, "bar": 2}
//	cp := shallowCopy(m) // cp is a new map containing the key-value pairs of m
//
//	empty := shallowCopy(map[string]any{}) // empty is an empty map
//
// Note: The function does not perform a deep copy, so any nested maps or
// reference types will be shared between the original and copied map.
func shallowCopy(m map[string]any) map[string]any {
	cp := make(map[string]any, len(m))

	if len(m) > 0 {
		for k, v := range m {
			cp[k] = v
		}
	}

	return cp
}

// merge performs a merge operation on two maps and returns a new map with the merged key-value pairs.
// It takes two parameters, m0 and m1, which are maps with string keys and values of any type.
// The function starts by making a shallow copy of the m0 map using the shallowCopy function.
// If m1 is not empty, it iterates over each key-value pair in m1 and adds them to the copied map.
// Finally, the function returns the copied map.
//
// Example:
//
//	m0 := map[string]any{"foo": 1, "bar": 2}
//	m1 := map[string]any{"baz": 3}
//	result := merge(m0, m1) // result is a new map containing {"foo": 1, "bar": 2, "baz": 3}
//
// Note: The function does not modify the original maps and returns a new map with the merged key-value pairs.
func merge(m0 map[string]any, m1 map[string]any) map[string]any {
	cp := shallowCopy(m0)

	if len(m1) > 0 {
		for k, v := range m1 {
			cp[k] = v
		}
	}

	return cp
}
