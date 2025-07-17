package logger

import (
	"fmt"
	"log"
	"runtime"
	"strings"
)

// LogLevel represents different log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	FATAL
)

// String returns string representation of log level
func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured logging with file and line information
type Logger struct {
	prefix string
}

// NewLogger creates a new logger instance
func NewLogger(prefix string) *Logger {
	return &Logger{
		prefix: prefix,
	}
}

// getCallerInfo returns file name and line number of the caller
func getCallerInfo(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "unknown", 0
	}

	// Get just the filename, not the full path
	parts := strings.Split(file, "/")
	filename := parts[len(parts)-1]

	return filename, line
}

// formatMessage formats the log message with caller info
func (l *Logger) formatMessage(level LogLevel, msg string, args ...interface{}) string {
	file, line := getCallerInfo(3) // Skip 3 levels: formatMessage -> log function -> caller

	formattedMsg := fmt.Sprintf(msg, args...)

	if l.prefix != "" {
		return fmt.Sprintf("[%s] %s %s:%d - %s", level.String(), l.prefix, file, line, formattedMsg)
	}

	return fmt.Sprintf("[%s] %s:%d - %s", level.String(), file, line, formattedMsg)
}

// Debug logs debug level message
func (l *Logger) Debug(msg string, args ...interface{}) {
	log.Println(l.formatMessage(DEBUG, msg, args...))
}

// Info logs info level message
func (l *Logger) Info(msg string, args ...interface{}) {
	log.Println(l.formatMessage(INFO, msg, args...))
}

// Warn logs warning level message
func (l *Logger) Warn(msg string, args ...interface{}) {
	log.Println(l.formatMessage(WARN, msg, args...))
}

// Error logs error level message
func (l *Logger) Error(msg string, args ...interface{}) {
	log.Println(l.formatMessage(ERROR, msg, args...))
}

// Fatal logs fatal level message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	log.Fatalln(l.formatMessage(FATAL, msg, args...))
}

// ErrorWithStack logs error with stack trace
func (l *Logger) ErrorWithStack(err error, msg string, args ...interface{}) {
	file, line := getCallerInfo(2)

	formattedMsg := fmt.Sprintf(msg, args...)

	logMsg := fmt.Sprintf("[ERROR] %s %s:%d - %s", l.prefix, file, line, formattedMsg)
	if err != nil {
		logMsg += fmt.Sprintf(" | Error: %v", err)
	}

	// Just log the error message without full stack trace
	log.Printf("%s", logMsg)
}

// Global logger instance
var defaultLogger = NewLogger("APP")

// Global logging functions for convenience

// Debug logs debug message using default logger
func Debug(msg string, args ...interface{}) {
	defaultLogger.Debug(msg, args...)
}

// Info logs info message using default logger
func Info(msg string, args ...interface{}) {
	defaultLogger.Info(msg, args...)
}

// Warn logs warning message using default logger
func Warn(msg string, args ...interface{}) {
	defaultLogger.Warn(msg, args...)
}

// Error logs error message using default logger
func Error(msg string, args ...interface{}) {
	defaultLogger.Error(msg, args...)
}

// Fatal logs fatal message using default logger
func Fatal(msg string, args ...interface{}) {
	defaultLogger.Fatal(msg, args...)
}

// ErrorWithStack logs error with stack trace using default logger
func ErrorWithStack(err error, msg string, args ...interface{}) {
	defaultLogger.ErrorWithStack(err, msg, args...)
}
