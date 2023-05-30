package logging

import "context"

const loggerKey = "logger"

// FromContext returns the Logger stored in the context, or a new one, if none is found.
func FromContext(ctx context.Context) Logger {
	if ctx.Value(loggerKey) == nil {
		return NewBasicLogger()
	}
	return ctx.Value(loggerKey).(Logger)
}

// ToContext adds the Logger to the context.
func ToContext(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, loggerKey, l)
}
