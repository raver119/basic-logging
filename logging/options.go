package logging

import "io"

type Option func(*Configuration)

// WithLogLevel sets the log level to use.
func WithLogLevel(level LogLevel) Option {
	return func(c *Configuration) {
		c.LogLevel = level
	}
}

// WithWriter sets the writer to use for logging. If nil, os.Stdout is used.
func WithWriter(writer io.Writer) Option {
	return func(c *Configuration) {
		c.Writer = writer
	}
}
