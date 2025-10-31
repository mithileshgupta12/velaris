package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"runtime/debug"
	"time"
)

type LogLevel int

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

type Field struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type Entry struct {
	Timestamp  time.Time      `json:"timestamp"`
	Level      string         `json:"level"`
	Message    string         `json:"message"`
	Fields     map[string]any `json:"fields,omitempty"`
	Stacktrace *string        `json:"stacktrace,omitempty"`
	Caller     *string        `json:"caller,omitempty"`
}

type Logger interface {
	Log(logLevel LogLevel, message string, fields []*Field)
}

type logger struct {
}

type testLogger struct {
}

func NewLogger() Logger {
	return &logger{}
}

func NewTestLogger() Logger {
	return &testLogger{}
}

func (l *logger) Log(logLevel LogLevel, message string, fields []*Field) {
	entry := &Entry{
		Timestamp: time.Now().UTC(),
		Level:     logLevel.String(),
		Message:   message,
	}

	if len(fields) > 0 {
		entry.Fields = make(map[string]any)
		for _, field := range fields {
			entry.Fields[field.Key] = field.Value
		}
	}

	if logLevel == ERROR || logLevel == FATAL {
		stacktrace := string(debug.Stack())

		_, file, line, _ := runtime.Caller(1)
		caller := fmt.Sprintf("%s:%d", file, line)

		entry.Stacktrace = &stacktrace
		entry.Caller = &caller
	}

	fmt.Println(l.formatJSON(entry))

	if logLevel == FATAL {
		os.Exit(1)
	}
}

func (l *logger) formatJSON(entry *Entry) string {
	data, err := json.Marshal(entry)
	if err != nil {
		return ""
	}

	return string(data)
}

func (l *testLogger) Log(logLevel LogLevel, message string, fields []*Field) {}
