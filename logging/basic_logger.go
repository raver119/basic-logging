package logging

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type BasicLogger struct {
	Configuration
	fields []Field
	lines  []Line

	manualEmit bool

	writer io.Writer
	lock   *sync.Mutex
}

func NewBasicLogger(options ...Option) *BasicLogger {
	configuration := DefaultConfiguration
	for _, option := range options {
		option(&configuration)
	}

	// safety check
	if configuration.Writer == nil {
		configuration.Writer = os.Stdout
	}

	return &BasicLogger{
		Configuration: configuration,
		lock:          &sync.Mutex{},
		writer:        configuration.Writer,
	}
}

func (l *BasicLogger) A(f ...Field) Logger {
	for _, field := range f {
		l.Add(field.Key, field.Value)
	}

	return l
}

// Add a field to the event
func (l *BasicLogger) Add(key string, value any) Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.fields = append(l.fields, Field{key, value})
	return l
}

func (l *BasicLogger) line(level LogLevel, format string, args ...any) Logger {
	l.lock.Lock()
	defer l.lock.Unlock()

	// if we use a manual emitter - we need to store the lines for later
	if l.manualEmit {
		l.lines = append(l.lines, Line{level, fmt.Sprintf(format, args...), time.Now().UTC()})
		return l
	}

	// otherwise we can just print the line

	return l
}

// Debugf logs a debug event
func (l *BasicLogger) Debugf(format string, args ...any) Logger {
	return l.line(LogLevelDebug, format, args...)
}

// Debugf logs a debug event
func (l *BasicLogger) Infof(format string, args ...any) Logger {
	return l.line(LogLevelInfo, format, args...)
}

// Warnf logs a warning event
func (l *BasicLogger) Warnf(format string, args ...any) Logger {
	return l.line(LogLevelWarn, format, args...)
}

// Errorf logs an error event
func (l *BasicLogger) Errorf(format string, args ...any) Logger {
	return l.line(LogLevelError, format, args...)
}

// Emitter returns a function that can be used to emit the logs.
// Once invoked, the logs will NOT be emitted automatically.
func (l *BasicLogger) Emitter() func() error {
	l.lock.Lock()
	defer l.lock.Unlock()

	l.manualEmit = true

	return func() error {
		l.lock.Lock()
		defer l.lock.Unlock()

		// if no lines above the log level - do not emit anything
		if len(l.lines) == 0 && len(l.fields) == 0 {
			return nil
		}

		output := map[string]interface{}{}

		// add fields first
		for _, field := range l.fields {
			output[field.Key] = field.Value
		}

		lines := []map[string]interface{}{}

		// process lines one by one
		for _, line := range l.lines {
			// skip all lines that are below the defined log level
			if line.Level < l.LogLevel {
				continue
			}

			lines = append(lines, line.ToMap())
		}

		output["lines"] = lines

		// timestamp goes the last
		output["ts"] = time.Now().UTC().UnixNano()

		b, err := json.Marshal(output)
		if err != nil {
			return fmt.Errorf("failed to marshal log entry: %w", err)
		}

		// write it away
		l.writer.Write(b)
		l.writer.Write([]byte("\n"))

		return nil
	}
}
