package log

// Name of fields, exhaustive list.
const (
	FieldListen   = "addr"     // Listened address and port.
	FieldServer   = "server"   // Name of the server.
	FieldError    = "error"    // Error causing the event.
	FieldSession  = "session"  // Current session.
	FieldCommand  = "command"  // Current command being processed.
	FieldResponse = "response" // Current response (to be) emited.
)

// Fields is used to define the content of an event with structured fields.
// It can be used as a shorthand for map[string]interface{}.
type Fields map[string]interface{}
