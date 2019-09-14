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
	"fmt"
	"os"
	"strings"
)

// Code is an alias for the type uint16
type Code uint16

// Response holds a 3 digit code and a messsage.
type Response struct {
	Code Code     `json:"code"`
	Msg  []string `json:"message"`
}

// IsError returns true if the response is an error.
func (e Response) IsError() bool {
	return strings.HasPrefix(e.String(), "5")
}

// IsSuccess returns true if the response is an success.
func (e Response) IsSuccess() bool {
	return !e.IsError()
}

func (e Response) String() string {
	var sb strings.Builder
	for i, msg := range e.Msg {
		if i == len(e.Msg)-1 {
			sb.WriteString(fmt.Sprintf("%3d %s", e.Code, msg))
		} else {
			sb.WriteString(fmt.Sprintf("%3d-%s\r\n", e.Code, msg))
		}
	}
	return sb.String()
}

// Resp is an alias for the type uint16
type Resp uint16

// Responses
const (
	Ready                 Resp = iota // First response
	Closing                           // Service closing
	Success                           // Requested action completed
	Abort                             // Requested action aborted
	Data                              // Ask for data input
	NotAvailable                      // Service is not available
	ShuttingDown                      // Service is shutting down
	SessionTimeout                    // Session timeout
	CommandUnrecognized               // Syntax error, command unrecognized
	ParameterSyntax                   // Syntax error in parameters or arguments
	CommandNotImplemented             // Command not implemented
	BadSequence                       // Bad sequence of commands
	NoValidRecipients                 // Transaction failed : no valid recipients
)

// SMTP reply codes as defined by RFC 5321, 4.2.3
const (
	CodeStatusHelp              Code = 211 // System status, or system help reply
	CodeHelp                    Code = 214 // Help message (Information on how to use the receiver or the meaning of a particular non-standard command; this reply is useful only to the human user)
	CodeReady                   Code = 220 // <domain> Service ready
	CodeClosing                 Code = 221 // <domain> Service closing transmission channel
	CodeSuccess                 Code = 250 // Requested mail action okay, completed
	CodeUserNotLocalTemp        Code = 251 // User not local; will forward to <forward-path>
	CodeCannotVerify            Code = 252 // Cannot VRFY user, but will accept message and attempt delivery
	CodeAskForData              Code = 354 // Start mail input; end with <CRLF>.<CRLF>
	CodeNotAvailable            Code = 421 // <domain> Service not available, closing transmission channel
	CodeMailboxUnavailableTemp  Code = 450 // Requested mail action not taken: mailbox unavailable (e.g., mailbox busy or temporarily blocked for policy reasons)
	CodeAbort                   Code = 451 // Requested action aborted: local error in processing
	CodeInsufficientStorageTemp Code = 452 // Requested action not taken: insufficient system storage
	CodeUnableAccomodateParam   Code = 455 // Server unable to accommodate parameters
	CodeCommandUnrecognized     Code = 500 // Syntax error, command unrecognized
	CodeParameterSyntax         Code = 501 // Syntax error in parameters or arguments
	CodeNotImplemented          Code = 502 // Command not implemented
	CodeBadSequence             Code = 503 // Bad sequence of commands
	CodeParameterNotImplemented Code = 504 // Command parameter not implemented
	CodeMailboxUnavailablePerm  Code = 550 // Requested action not taken: mailbox unavailable (e.g., mailbox not found, no access, or command rejected for policy reasons)
	CodeUserNotLocalPerm        Code = 551 // User not local; please try <forward-path>
	CodeInsufficientStoragePerm Code = 552 // Requested mail action aborted: exceeded storage allocation
	CodeMailboxNotAllowed       Code = 553 // Requested action not taken: mailbox name not allowed (e.g., mailbox syntax incorrect)
	CodeTransactionFailed       Code = 554 // Transaction failed
	CodeMailFromRcptToParam     Code = 555 // MAIL FROM/RCPT TO parameters not recognized or not implemented
)

// Responses returned by the SMTP server
var Responses = map[Resp]Response{
	Ready:                 Response{CodeReady, []string{"<domain> Service ready"}},
	Closing:               Response{CodeClosing, []string{"<domain> Service closing transmission channel"}},
	Success:               Response{CodeSuccess, []string{"OK"}},
	Data:                  Response{CodeAskForData, []string{"Start mail input; end with <CRLF>.<CRLF>"}},
	NotAvailable:          Response{CodeNotAvailable, []string{"<domain> Service not available, closing transmission channel"}},
	ShuttingDown:          Response{CodeNotAvailable, []string{"<domain> Service shutting down and closing transmission channel"}},
	SessionTimeout:        Response{CodeNotAvailable, []string{"Your session timed out due to inactivity"}},
	Abort:                 Response{CodeAbort, []string{"Requested action aborted: error in processing"}},
	CommandUnrecognized:   Response{CodeCommandUnrecognized, []string{"Syntax error, command unrecognized"}},
	ParameterSyntax:       Response{CodeParameterSyntax, []string{"Syntax error in parameters or arguments"}},
	CommandNotImplemented: Response{CodeNotImplemented, []string{"Command not implemented"}},
	BadSequence:           Response{CodeBadSequence, []string{"Bad sequence of commands"}},
	NoValidRecipients:     Response{CodeTransactionFailed, []string{"No valid recipients"}},
}

var hostname string

// SetReply set reply text of given code.
func SetReply(r Resp, s ...string) {
	response := Responses[r]
	response.Msg = s
	for i, msg := range response.Msg {
		response.Msg[i] = strings.ReplaceAll(msg, "<domain>", hostname)
	}
	Responses[r] = response
}

func r(r Resp) *Response {
	response := Responses[r]
	return &response
}

func init() {
	var err error
	if hostname, err = os.Hostname(); err != nil {
		hostname = "localhost"
	}
	for code, text := range Responses {
		SetReply(code, text.Msg...)
	}
}
