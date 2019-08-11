// Package smtpd contains source code of the SMTP server of Mailmock
//
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
//
// Linking this library statically or dynamically with other modules is
// making a combined work based on this library.  Thus, the terms and
// conditions of the GNU General Public License cover the whole
// combination.
//
// As a special exception, the copyright holders of this library give you
// permission to link this library with independent modules to produce an
// executable, regardless of the license terms of these independent
// modules, and to copy and distribute the resulting executable under
// terms of your choice, provided that you also meet, for each linked
// independent module, the terms and conditions of the license of that
// module.  An independent module is a module which is not derived from
// or based on this library.  If you modify this library, you may extend
// this exception to your version of the library, but you are not
// obligated to do so.  If you do not wish to do so, delete this
// exception statement from your version.
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
	"VRFY": {1, true, []string{""}},
}

// ParseCommand parses a SMTP command, returns appropriate response if the command is malformed
// If the command is well formed, returned response is nil
func ParseCommand(input string) (*Command, *Response) {
	elmts := strings.Split(input, " ")
	name := strings.ToUpper(strings.TrimSpace(elmts[0]))
	desc, ok := listOfValidCommands[name]
	if !ok {
		return nil, &Response{500, "Syntax error, command unrecognized"}
	}

	if desc.numberOfArgument != len(desc.argumentNames) {
		panic("Coding Error: wrong cmdDescription for " + name)
	}

	elmts = elmts[1:]

	if len(elmts) < desc.numberOfArgument {
		return nil, &Response{501, "Syntax error in parameters or arguments"}
	}

	if len(elmts) > desc.numberOfArgument && desc.isStrict {
		return nil, &Response{501, "Syntax error in parameters or arguments"}
	}

	command := &Command{FullCmd: input, Name: name, PositionalArgs: []string{}, NamedArgs: map[string]string{}}

	for i, arg := range elmts {
		argPos := i - len(command.PositionalArgs)
		if argPos >= len(desc.argumentNames) {
			break
		}
		argName := desc.argumentNames[argPos]
		if argName != "" {
			if strings.Count(arg, ":") != 1 {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
			argSplit := strings.Split(arg, ":")
			if strings.ToUpper(argSplit[0]) != argName {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
			command.NamedArgs[argName] = strings.TrimSpace(argSplit[1])
			if command.NamedArgs[argName] == "" {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
		} else {
			if arg == "" {
				return nil, &Response{501, "Syntax error in parameters or arguments"}
			}
			command.PositionalArgs = append(command.PositionalArgs, strings.TrimSpace(arg))
		}
	}

	return command, nil
}
