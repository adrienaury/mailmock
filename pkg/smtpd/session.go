package smtpd

import (
	"log"
	"net/textproto"
)

// SessionState is the state of a Session
type SessionState string

// Session States
const (
	SSInitiated SessionState = "initiated"
	SSReady     SessionState = "ready"
	SSBusy      SessionState = "busy"
	SSClosed    SessionState = "closed"
)

// TransactionHandler will be called each time a transaction reach TSCompleted or TSAborted status
type TransactionHandler func(*Transaction)

// Session :
type Session struct {
	state  SessionState
	client string
	tr     *Transaction
	conn   *textproto.Conn
	th     *TransactionHandler
}

// NewSession return a new Session
func NewSession(c *textproto.Conn, th *TransactionHandler) *Session {
	s := &Session{state: SSInitiated, conn: c, th: th}
	log.Printf("[%p] New session %v ", s, s)
	return s
}

// Serve will reponds to any request until a QUIT command is received
func (s *Session) Serve() {
	s.conn.PrintfLine("%v", Response{220, "Service ready"})

	for {
		input, err := s.conn.ReadLine()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("[%p] %9s %15s => %v", s, s.state, s.client, input)
		r := s.receive(input)
		log.Printf("[%p] %9s %15s <= %v", s, s.state, s.client, r)
		err = s.conn.PrintfLine("%v", r)
		if err != nil {
			log.Fatal(err)
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
		res = s.hello(cmd)
	case "EHLO":
		res = s.hello(cmd)
	case "MAIL":
		res = s.mail(cmd)
	case "RCPT":
		res = s.rcpt(cmd)
	case "DATA":
		res = s.data(cmd)
	case "NOOP":
		res = s.noop(cmd)
	case "RSET":
		res = s.reset(cmd)
	case "QUIT":
		res = s.quit(cmd)
	case "VRFY":
		res = s.verify(cmd)
	default:
		log.Fatal("Coding Error")
	}
	return res
}

func (s *Session) hello(cmd *Command) *Response {
	s.client = cmd.PositionalArgs[0]
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
		// TODO
	}
	return res
}

func (s *Session) rcpt(cmd *Command) *Response {
	if s.state != SSBusy {
		return &Response{503, "Bad sequence of commands"}
	}
	res, err := s.tr.Process(cmd)
	if err != nil {
		// TODO
	}
	return res
}

func (s *Session) data(cmd *Command) *Response {
	if s.state != SSBusy {
		return &Response{503, "Bad sequence of commands"}
	}

	res, err := s.tr.Process(cmd)
	if err != nil {
		// TODO
	}

	s.conn.PrintfLine("%v", res)
	data, err := s.conn.ReadDotLines()
	if err != nil {
		// TODO
	}

	res, err = s.tr.Data(data)
	if err != nil {
		// TODO
	}

	s.state = SSReady
	return res
}

func (s *Session) verify(cmd *Command) *Response {
	return &Response{502, "Command not implemented"}
}

func (s *Session) noop(cmd *Command) *Response {
	return &Response{250, "OK"}
}

func (s *Session) reset(cmd *Command) *Response {
	err := s.tr.Abort()
	if err != nil {
		// TODO
	}

	if s.client != "" {
		s.state = SSReady
	} else {
		s.state = SSInitiated
	}

	return &Response{250, "OK"}
}

func (s *Session) quit(cmd *Command) *Response {
	s.state = SSClosed
	s.tr.Abort()
	return &Response{Code: 221, Msg: "Service closing transmission channel"}
}

func (s *Session) handleTransaction() {
	if (*s.th) != nil && s.tr != nil {
		go (*s.th)(s.tr)
	}
	s.tr = nil
}
