package logging

import (
	"io"
	"os"
)

type Configuration struct {
	LogLevel LogLevel
	Writer   io.Writer
}

var DefaultConfiguration = Configuration{
	LogLevel: LogLevelInfo,
	Writer:   os.Stdout,
}
