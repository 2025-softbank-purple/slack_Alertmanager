package client

import (
	"github.com/sirupsen/logrus"
)

// Logger provides structured logging functionality
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return &Logger{logger: logger}
}

// Info logs an info message
func (l *Logger) Info(message string) error {
	l.logger.Info(message)
	return nil
}

// Error logs an error message
func (l *Logger) Error(message string) error {
	l.logger.Error(message)
	return nil
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.logger.WithFields(fields)
}

