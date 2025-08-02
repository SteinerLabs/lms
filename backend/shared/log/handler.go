package log

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"runtime"
	"strconv"
	"sync"
)

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
