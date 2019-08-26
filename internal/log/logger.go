// Copyright (C) 2019  Adrien Aury
//
// This file is part of Mailmock.
//
// Mailmock is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mailmock is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Mailmock.  If not, see <https://www.gnu.org/licenses/>.

// Package log contains the Logger interface and everything related.
package log

import (
	"github.com/goph/logur"
)

// Logger is the fundamental interface for all log operations.
type Logger interface {
	Trace(msg string, fields ...map[string]interface{})
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})

	// WithFields annotates a logger with some context and return it as a new instance.
	WithFields(fields map[string]interface{}) Logger
}

// LoggerNoop is a Logger implementation that does nothing.
type LoggerNoop struct{}

// Trace for LoggerNoop does nothing.
func (LoggerNoop) Trace(msg string, fields ...map[string]interface{}) {}

// Debug for LoggerNoop does nothing.
func (LoggerNoop) Debug(msg string, fields ...map[string]interface{}) {}

// Info for LoggerNoop does nothing.
func (LoggerNoop) Info(msg string, fields ...map[string]interface{}) {}

// Warn for LoggerNoop does nothing.
func (LoggerNoop) Warn(msg string, fields ...map[string]interface{}) {}

// Error for LoggerNoop does nothing.
func (LoggerNoop) Error(msg string, fields ...map[string]interface{}) {}

// WithFields for LoggerNoop does nothing.
func (l LoggerNoop) WithFields(fields map[string]interface{}) Logger { return l }

// LoggerAdapter is a Logger implementation that wrap a Logur logger.
type LoggerAdapter struct {
	logger logur.Logger
}

// Trace logs a trace event.
func (l *LoggerAdapter) Trace(msg string, fields ...map[string]interface{}) {
	l.logger.Trace(msg, fields...)
}

// Debug logs a debug event.
func (l *LoggerAdapter) Debug(msg string, fields ...map[string]interface{}) {
	l.logger.Debug(msg, fields...)
}

// Info logs an info event.
func (l *LoggerAdapter) Info(msg string, fields ...map[string]interface{}) {
	l.logger.Info(msg, fields...)
}

// Warn logs a warning event.
func (l *LoggerAdapter) Warn(msg string, fields ...map[string]interface{}) {
	l.logger.Warn(msg, fields...)
}

// Error logs an error event.
func (l *LoggerAdapter) Error(msg string, fields ...map[string]interface{}) {
	l.logger.Error(msg, fields...)
}

// WithFields annotates a logger with some context and it as a new instance.
func (l *LoggerAdapter) WithFields(fields map[string]interface{}) Logger {
	return &LoggerAdapter{logger: logur.WithFields(l.logger, fields)}
}

// NewLoggerAdapter returns a new Logger instance.
func NewLoggerAdapter(logger logur.Logger) Logger {
	return &LoggerAdapter{
		logger: logger,
	}
}

// DefaultLogger is the logger user by this package.
var DefaultLogger Logger = LoggerNoop{}
