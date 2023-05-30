package logging

import "time"

type Line struct {
	Level   LogLevel
	Message string
	Time    time.Time
}

func (l *Line) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"level":   l.Level,
		"message": l.Message,
		"time":    l.Time,
	}
}
