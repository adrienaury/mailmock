package smtpd

import (
	"strings"
)

// Command is a parsed SMTP command
type Command struct {
	FullCmd        string
	Name           string
	PositionalArgs []string
	NamedArgs      map[string]string
}

type cmdDescription struct {
	numberOfArgument int      // expected number of argument
	isStrict         bool     // true if additional arguments are prohibited
	argumentNames    []string // names of arguments in order (prefixes)
}

var listOfValidCommands = map[string]cmdDescription{
	"HELO": {1, true, []string{""}},
	"EHLO": {1, true, []string{""}},
	"MAIL": {1, true, []string{"FROM"}},
	"RCPT": {1, true, []string{"TO"}},
	"DATA": {0, true, []string{}},
	"NOOP": {0, false, []string{}},
	"RSET": {0, true, []string{}},
	"QUIT": {0, true, []string{}},
	"VRFY": {0, true, []string{}},
}

// ParseCommand parses a SMTP command, returns appropriate response if the command is malformed
// If the command is well formed, returned response is nil
func ParseCommand(cmd string) (*Command, *Response) {
	elmts := strings.Split(cmd, " ")
	name := strings.ToUpper(strings.TrimSpace(elmts[0]))
	desc, ok := listOfValidCommands[name]
	if !ok {
		return nil, &Response{500, "Syntax error, command unrecognized"}
	}

	elmts = elmts[1:]

	if len(elmts) < desc.numberOfArgument {
		return nil, &Response{501, "Syntax error in parameters or arguments"}
	}

	if len(elmts) > desc.numberOfArgument && desc.isStrict {
		return nil, &Response{501, "Syntax error in parameters or arguments"}
	}

	command := &Command{FullCmd: cmd, Name: name, PositionalArgs: []string{}, NamedArgs: map[string]string{}}

	for i, arg := range elmts {
		argName := desc.argumentNames[i]
		if argName != "" {
			if strings.Count(arg, ":") != 1 {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
			argSplit := strings.Split(arg, ":")
			if strings.ToUpper(argSplit[0]) != argName {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
			command.NamedArgs[argName] = argSplit[1]
		} else {
			command.PositionalArgs[i] = arg
		}
	}

	return command, nil
}
