package log

import (
	"io"
	"log/slog"
	"os"
)

var defaultConfig = &Options{
	Output:     os.Stdout,
	AddSource:  true,
	Level:      defaultLogLevel,
	TimeFormat: "",
	Prefix:     "",
	Json:       false,
}

// Option uses the functional options patter to configure a Logger
type Option func(*Options)

type Options struct {
	Output     io.Writer
	AddSource  bool
	Level      slog.Leveler
	TimeFormat string
	Prefix     string
	Json       bool
}

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
