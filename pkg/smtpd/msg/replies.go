package msg

var (
	// GreetingBanner is the first reply after a client opens a new connection.
	GreetingBanner = "Service ready"

	// Success is the default reply after successful processing.
	Success = "OK"

	// AskForData is replied to a DATA command before receiving mail content.
	AskForData = "Start mail input; end with <CRLF>.<CRLF>"

	// Goodbye is the last reply sent before closing the transmission channel.
	Goodbye = "Service closing transmission channel"

	// ServiceNotAvailable is sent when the server cannot accept new connection.
	ServiceNotAvailable = "Service not available, closing transmission channel"

	// RequestedActionAborted is sent when an error occured processing the request.
	RequestedActionAborted = "Requested action aborted: error in processing"

	// BadSequence is sent when the last command was correct, but not expected.
	BadSequence = "Bad sequence of commands"

	// NoValidRecipients is sent when no valid recipients were issued.
	NoValidRecipients = "No valid recipients"

	// NotImplemented is sent when the command is recognized, but not implemented by the server.
	NotImplemented = "Command not implemented"

	// NotRecognized is sent when the command is not recognized.
	NotRecognized = "Syntax error, command unrecognized"

	// ParameterError is sent when the command is recognized, but arguments or parameters are invalid.
	ParameterError = "Syntax error in parameters or arguments"
)
