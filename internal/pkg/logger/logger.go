package logger

import (
	"encoding/json"
	"fmt"
	"time"
)

type LogLevel int

const (
	colorReset = "\033[0m"
	colorGreen = "\033[32m"
)

const (
	INFO LogLevel = iota
	DEBUG
	ERROR
	FATAL
)

func (l LogLevel) String() string {
	switch l {
	case INFO:
		return "INFO"
	case DEBUG:
		return "DEBUG"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

type Format int

const (
	FormatJSON Format = iota
	FormatHuman
)

type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
}

type Logger interface {
	Log(format Format, logLevel LogLevel, message string)
}

type logger struct{}

func NewLogger() Logger {
	return &logger{}
}

func (l *logger) Log(format Format, logLevel LogLevel, message string) {
	entry := &Entry{
		Timestamp: time.Now().UTC(),
		Level:     logLevel.String(),
		Message:   message,
	}

	switch format {
	case FormatJSON:
		fmt.Println(l.formatJSON(entry))
	case FormatHuman:
		fmt.Println(l.formatHuman(entry))
	}

}

func (l *logger) formatJSON(entry *Entry) string {
	data, err := json.Marshal(entry)
	if err != nil {
		return ""
	}

	return string(data)
}

func (l *logger) formatHuman(entry *Entry) string {
	return fmt.Sprintf("[%s] %s%s%s: %s",
		entry.Timestamp,
		colorGreen,
		entry.Level,
		colorReset,
		entry.Message,
	)
}
