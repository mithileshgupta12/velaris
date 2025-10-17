package logger

import (
	"fmt"
	"time"
)

type LogLevel int

const (
	INFO LogLevel = iota
	DEBUG
	ERROR
)

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (l *Logger) Log(logLevel LogLevel, message string) {
	fmt.Printf("[%s] %s: %s\n", time.Now().UTC(), logLevel, message)
}
