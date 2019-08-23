package log

// Logger is the fundamental interface for all log operations.
type Logger interface {
	Trace(msg string, fields ...map[string]interface{})
	Debug(msg string, fields ...map[string]interface{})
	Info(msg string, fields ...map[string]interface{})
	Warn(msg string, fields ...map[string]interface{})
	Error(msg string, fields ...map[string]interface{})
}

// NoopLogger is a Logger implementation that does nothing.
type NoopLogger struct{}

// Trace for NoopLogger does nothing.
func (NoopLogger) Trace(msg string, fields ...map[string]interface{}) {}

// Debug for NoopLogger does nothing.
func (NoopLogger) Debug(msg string, fields ...map[string]interface{}) {}

// Info for NoopLogger does nothing.
func (NoopLogger) Info(msg string, fields ...map[string]interface{}) {}

// Warn for NoopLogger does nothing.
func (NoopLogger) Warn(msg string, fields ...map[string]interface{}) {}

// Error for NoopLogger does nothing.
func (NoopLogger) Error(msg string, fields ...map[string]interface{}) {}

// DefaultLogger is the logger user by this package.
var DefaultLogger Logger = NoopLogger{}
