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
	"io"
	"net"
	"net/textproto"
	"time"

	"github.com/adrienaury/mailmock/internal/log"
)

// SessionState is the state of a Session.
type SessionState string

// Session States
const (
	SSInitiated SessionState = "initiated"
	SSReady     SessionState = "ready"
	SSBusy      SessionState = "busy"
	SSClosed    SessionState = "closed"
)

// TransactionHandler will be called each time a transaction reach TSCompleted or TSAborted status.
type TransactionHandler func(*Transaction)

// Session represents a SMTP session of a client.
type Session struct {
	State    SessionState `json:"state"`
	Client   string       `json:"client"`
	Tr       *Transaction `json:"transaction"`
	conn     *textproto.Conn
	th       *TransactionHandler
	logger   log.Logger
	tcpConn  *net.TCPConn
	mustStop bool
}

// NewSession return a new Session.
func NewSession(c *textproto.Conn, th *TransactionHandler, logger log.Logger) *Session {
	s := &Session{State: SSInitiated, conn: c, th: th, logger: nil}
	if logger == nil {
		logger = log.DefaultLogger
	}
	l := logger.WithFields(log.Fields{
		log.FieldSession: s,
	})
	s.logger = l
	s.logger.Info("Initiated new session")
	return s
}

// Serve will reponds to any request until a QUIT command is received or connection is broken.
func (s *Session) Serve(stop <-chan struct{}) {
	if s.State == SSClosed {
		s.logger.Warn("Cannot serve a closed session")
		if err := s.conn.PrintfLine("%v", r(CodeNotAvailable)); err != nil {
			s.logger.Error("Failed to send response to client", log.Fields{log.FieldError: err, log.FieldResponse: r(CodeNotAvailable)})
		}
		return
	}

	if err := s.conn.PrintfLine("%v", r(CodeReady)); err != nil {
		s.logger.Error("Failed to send greeting message, quitting session", log.Fields{log.FieldError: err, log.FieldResponse: r(CodeReady)})
		s.quit()
		return
	}

	shutdown := make(chan struct{})
	defer close(shutdown)

	go func() {
		// Block until either stop or shutdown signal
		select {
		case <-stop:
			s.mustStop = true
			s.logger.Warn("Server must stop, session will timeout in 30 seconds (at most)")
			<-time.After(30 * time.Second)
			if s.tcpConn != nil {
				_ = s.tcpConn.SetReadDeadline(time.Now())
			}
		case <-shutdown:
		}
		<-time.After(5 * time.Second)
		s.conn.Close()
	}()

	s.serveLoop(stop)
}

func (s *Session) serveLoop(stop <-chan struct{}) {
	for s.State != SSClosed {
		var res *Response

		if s.tcpConn != nil {
			// SMTP server SHOULD have a timeout of at least 5 minutes while it
			// is awaiting the next command from the sender (RFC 5321 4.5.3.2.7.)
			if err := s.tcpConn.SetReadDeadline(time.Now().Add(time.Minute * 5)); err != nil {
				s.logger.Error("SetDeadline on SMTP session failed", log.Fields{log.FieldError: err})
			}
		}

		if input, err := s.conn.ReadLine(); err == io.EOF || err == io.ErrClosedPipe {
			s.logger.Error("Lost client connection, quitting", log.Fields{log.FieldError: err})
			res = s.quit()
		} else if errop, ok := err.(net.Error); ok && errop.Timeout() {
			if s.mustStop {
				s.logger.Warn("Session interrupted because server is shutting down")
			} else {
				s.logger.Warn("Session timed out")
			}
			if err := s.conn.PrintfLine("%v", r(CodeNotAvailable)); err != nil {
				s.logger.Error("Failed to send response to client", log.Fields{log.FieldError: err, log.FieldResponse: r(CodeNotAvailable)})
			}
			return
		} else if err != nil {
			s.logger.Error("Network error, requested action cannot be processed", log.Fields{log.FieldError: err})
			res = r(CodeAbort)
		} else {
			s.logger.Debug("Received command", log.Fields{log.FieldCommand: input})
			res = s.receive(input)
			if res.IsError() {
				s.logger.Warn("Processed command", log.Fields{log.FieldCommand: input, log.FieldResponse: res})
			} else {
				s.logger.Info("Processed command", log.Fields{log.FieldCommand: input, log.FieldResponse: res})
			}
		}

		select {
		case <-stop:
			// We need to shutdown
			s.logger.Warn("Session interrupted because server is shutting down")
			if err := s.conn.PrintfLine("%v", r(CodeNotAvailable)); err != nil {
				s.logger.Error("Failed to send response to client", log.Fields{log.FieldError: err, log.FieldResponse: r(CodeNotAvailable)})
			}
			return
		default:
		}

		if err := s.conn.PrintfLine("%v", res); err != nil {
			s.logger.Error("Network error, failed to send response, quitting", log.Fields{log.FieldError: err, log.FieldResponse: r(CodeNotAvailable)})
			s.quit()
			return
		}
		s.handleTransaction()
	}
}

func (s *Session) receive(input string) (res *Response) {
	cmd, res := ParseCommand(input)
	if res != nil {
		return res
	}
	switch cmd.Name {
	case "HELO":
		res = s.hello(cmd.PositionalArgs[0])
	case "EHLO":
		res = s.hello(cmd.PositionalArgs[0])
	case "MAIL":
		res = s.mail(cmd)
	case "RCPT":
		res = s.rcpt(cmd)
	case "DATA":
		res = s.data(cmd)
	case "NOOP":
		res = s.noop()
	case "RSET":
		res = s.reset()
	case "QUIT":
		res = s.quit()
	case "VRFY":
		res = s.verify(cmd.PositionalArgs[0])
	default:
		s.logger.Error("Coding error, this should not happen")
	}
	return res
}

func (s *Session) hello(client string) *Response {
	s.Client = client
	s.State = SSReady
	return r(CodeSuccess)
}

func (s *Session) mail(cmd *Command) *Response {
	if s.State != SSReady {
		return r(CodeBadSequence)
	}
	s.Tr = NewTransaction()
	s.logger.Debug("Started transaction")
	res, err := s.Tr.Process(cmd)
	if err != nil {
		return r(CodeAbort)
	}
	s.State = SSBusy
	return res
}

func (s *Session) rcpt(cmd *Command) *Response {
	if s.State != SSBusy {
		return r(CodeBadSequence)
	}
	res, err := s.Tr.Process(cmd)
	if err != nil {
		return r(CodeAbort)
	}
	return res
}

func (s *Session) data(cmd *Command) *Response {
	if s.State != SSBusy {
		return r(CodeBadSequence)
	}
	if len(s.Tr.Mail.Envelope.Recipients) == 0 {
		return r(CodeTransactionFailed)
	}

	res, err := s.Tr.Process(cmd)
	if err != nil {
		return r(CodeAbort)
	}

	if err = s.conn.PrintfLine("%v", res); err != nil {
		s.logger.Error("Failed to send response to client", log.Fields{log.FieldError: err, log.FieldResponse: res})
		return r(CodeAbort)
	}
	data, err := s.conn.ReadDotLines()
	if err != nil {
		return r(CodeAbort)
	}

	res, err = s.Tr.Data(data)
	if err != nil {
		return r(CodeAbort)
	}

	s.State = SSReady
	return res
}

func (s *Session) verify(address string) *Response {
	return r(CodeNotImplemented)
}

func (s *Session) noop() *Response {
	return r(CodeSuccess)
}

func (s *Session) reset() *Response {
	err := s.Tr.Abort()
	if err != nil {
		return r(CodeAbort)
	}

	if s.Client != "" {
		s.State = SSReady
	} else {
		s.State = SSInitiated
	}

	return r(CodeSuccess)
}

func (s *Session) quit() *Response {
	s.State = SSClosed
	_ = s.Tr.Abort()
	return r(CodeClosing)
}

func (s *Session) handleTransaction() {
	if s.Tr != nil && (s.Tr.State == TSCompleted || s.Tr.State == TSAborted) {
		s.logger.Debug("Ended transaction")
		if s.th != nil && (*s.th) != nil {
			(*s.th)(s.Tr)
		}
		s.Tr = nil
	}
}

func (s *Session) String() string {
	return fmt.Sprintf("%p[%v]", s, s.State)
}
