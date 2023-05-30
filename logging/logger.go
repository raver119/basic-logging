package logging

type Logger interface {
	// Methods for logging at different levels
	Debugf(format string, args ...any) Logger
	Infof(format string, args ...any) Logger
	Warnf(format string, args ...any) Logger
	Errorf(format string, args ...any) Logger

	// Add a field to the logger
	Add(key string, value any) Logger
}

func New(options ...Option) *BasicLogger {
	return NewBasicLogger(options...)
}
