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
	"io"
	"log"
	"net/textproto"
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
	state  SessionState
	client string
	tr     *Transaction
	conn   *textproto.Conn
	th     *TransactionHandler
}

// NewSession return a new Session.
func NewSession(c *textproto.Conn, th *TransactionHandler) *Session {
	s := &Session{state: SSInitiated, conn: c, th: th}
	log.Printf("[%p] New session initiated", s)
	return s
}

// Serve will reponds to any request until a QUIT command is received or connection is broken.
func (s *Session) Serve() {
	if s.state == SSClosed {
		s.conn.PrintfLine("%v", Response{421, "Service not available, closing transmission channel"})
		return
	}

	if err := s.conn.PrintfLine("%v", Response{220, "Service ready"}); err != nil {
		s.quit()
		log.Printf("[%p] %9s %15s == connection error : %v", s, s.state, s.client, err)
		return
	}

	for {
		var r *Response

		if input, err := s.conn.ReadLine(); err == io.EOF || err == io.ErrClosedPipe {
			r = s.quit()
			log.Printf("[%p] %9s %15s == client closed connection", s, s.state, s.client)
		} else if err != nil {
			log.Printf("[%p] %9s %15s == connection error : %v", s, s.state, s.client, err)
			r = &Response{451, "Requested action aborted: error in processing"}
		} else {
			log.Printf("[%p] %9s %15s => %v", s, s.state, s.client, input)
			r = s.receive(input)
		}

		log.Printf("[%p] %9s %15s <= %v", s, s.state, s.client, r)
		if err := s.conn.PrintfLine("%v", r); err != nil {
			s.quit()
			log.Printf("[%p] %9s %15s == connection error : %v", s, s.state, s.client, err)
			return
		}
		if s.tr != nil && (s.tr.State == TSCompleted || s.tr.State == TSAborted) {
			s.handleTransaction()
		}
		if s.state == SSClosed {
			break
		}
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
		log.Fatal("Coding Error")
	}
	return res
}

func (s *Session) hello(client string) *Response {
	s.client = client
	s.state = SSReady
	return &Response{250, "OK"}
}

func (s *Session) mail(cmd *Command) *Response {
	if s.state != SSReady {
		return &Response{503, "Bad sequence of commands"}
	}
	s.tr = NewTransaction()
	res, err := s.tr.Process(cmd)
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}
	s.state = SSBusy
	return res
}

func (s *Session) rcpt(cmd *Command) *Response {
	if s.state != SSBusy {
		return &Response{503, "Bad sequence of commands"}
	}
	res, err := s.tr.Process(cmd)
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}
	return res
}

func (s *Session) data(cmd *Command) *Response {
	if s.state != SSBusy {
		return &Response{503, "Bad sequence of commands"}
	}
	if len(s.tr.Mail.Envelope.Recipients) == 0 {
		return &Response{554, "No valid recipients"}
	}

	res, err := s.tr.Process(cmd)
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}

	s.conn.PrintfLine("%v", res)
	data, err := s.conn.ReadDotLines()
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}

	res, err = s.tr.Data(data)
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}

	s.state = SSReady
	return res
}

func (s *Session) verify(address string) *Response {
	return &Response{502, "Command not implemented"}
}

func (s *Session) noop() *Response {
	return &Response{250, "OK"}
}

func (s *Session) reset() *Response {
	err := s.tr.Abort()
	if err != nil {
		return &Response{451, "Requested action aborted: error in processing"}
	}

	if s.client != "" {
		s.state = SSReady
	} else {
		s.state = SSInitiated
	}

	return &Response{250, "OK"}
}

func (s *Session) quit() *Response {
	s.state = SSClosed
	s.tr.Abort()
	return &Response{Code: 221, Msg: "Service closing transmission channel"}
}

func (s *Session) handleTransaction() {
	if s.th != nil && (*s.th) != nil && s.tr != nil {
		go (*s.th)(s.tr)
	}
	s.tr = nil
}
