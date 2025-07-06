package logger

import (
	"log"
	"os"
)

// Logger represents the logger interface
type Logger interface {
	Info(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Debug(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
}

// logger implements the Logger interface
type logger struct {
	infoLogger  *log.Logger
	errorLogger *log.Logger
	debugLogger *log.Logger
	warnLogger  *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() Logger {
	return &logger{
		infoLogger:  log.New(os.Stdout, "[INFO] ", log.LstdFlags),
		errorLogger: log.New(os.Stderr, "[ERROR] ", log.LstdFlags),
		debugLogger: log.New(os.Stdout, "[DEBUG] ", log.LstdFlags),
		warnLogger:  log.New(os.Stdout, "[WARN] ", log.LstdFlags),
	}
}

// Info logs info messages
func (l *logger) Info(msg string, args ...interface{}) {
	l.infoLogger.Printf(msg, args...)
}

// Error logs error messages
func (l *logger) Error(msg string, args ...interface{}) {
	l.errorLogger.Printf(msg, args...)
}

// Debug logs debug messages
func (l *logger) Debug(msg string, args ...interface{}) {
	l.debugLogger.Printf(msg, args...)
}

// Warn logs warning messages
func (l *logger) Warn(msg string, args ...interface{}) {
	l.warnLogger.Printf(msg, args...)
}
