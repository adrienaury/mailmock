package smtpd_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/textproto"
	"strings"
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

type MockConn struct {
	snd *bytes.Buffer
	rcv *bytes.Buffer
}

func (c *MockConn) Read(p []byte) (n int, err error) {
	return c.snd.Read(p)
}

func (c *MockConn) Write(p []byte) (n int, err error) {
	return c.rcv.Write(p)
}

func (c *MockConn) Close() error {
	return nil
}

type EOFConn struct{}

func (c *EOFConn) Read(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (c *EOFConn) Write(p []byte) (n int, err error) {
	return 0, io.EOF
}

func (c *EOFConn) Close() error {
	return nil
}

func TestSessionNominal(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"MAIL FROM:<sender@example.com>",
			"RCPT TO:<recipient@example.com>",
			"DATA",
			"Subject: Test",
			"",
			"This is a test",
			".",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"250 OK",
			"250 OK",
			"354 Start mail input; end with <CRLF>.<CRLF>",
			"250 OK",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionReset1(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"MAIL FROM:<sender@example.com>",
			"RCPT TO:<recipient@example.com>",
			"RSET",
			"MAIL FROM:<sender@example.com>",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"250 OK",
			"250 OK",
			"250 OK",
			"250 OK",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionReset2(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"RSET",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionNoop(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"NOOP",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionVerify(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"VRFY test",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"502 Command not implemented",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionHeloHelp(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"HELP",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"500 Syntax error, command unrecognized",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionEhloHelp(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"EHLO localhost",
			"HELP",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK (extended)",
			"214 ",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionInvalidCommand(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"FAKE test",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"500 Syntax error, command unrecognized",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionBadSequence1(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"RCPT TO:<recipient@example.com>",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"503 Bad sequence of commands",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionBadSequence2(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"MAIL FROM:<sender@example.com>",
			"MAIL FROM:<sender@example.com>",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"250 OK",
			"503 Bad sequence of commands",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestSessionBadSequence3(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"DATA",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"503 Bad sequence of commands",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestNoValidRecipients(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"MAIL FROM:test",
			"DATA",
			"QUIT",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"250 OK",
			"554 No valid recipients",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	test(t, snd, rcv)
}

func TestClosedConnection(t *testing.T) {
	var (
		snd string = strings.Join([]string{
			"HELO localhost",
			"QUIT",
			"HELO localhost",
		}, "\r\n")
		rcv string = strings.Join([]string{
			"220 Service ready",
			"250 OK",
			"221 Service closing transmission channel",
			"",
		}, "\r\n")
	)
	session, rwc := test(t, snd, rcv)

	var (
		snd2 string = strings.Join([]string{
			"HELO localhost",
		}, "\r\n")
		rcv2 string = strings.Join([]string{
			"421 Service not available, closing transmission channel",
			"",
		}, "\r\n")
	)

	test2(t, snd2, rcv2, session, rwc.rcv)
}

func TestEOFConnection(t *testing.T) {
	rwc := &EOFConn{}

	c := textproto.NewConn(rwc)
	assert.NotNil(t, c, "")

	s := smtpd.NewSession(c, nil, nil)
	assert.NotNil(t, s, "")

	test2(t, strings.Join([]string{
		"HELO localhost",
	}, "\r\n"), strings.Join([]string{
		"",
	}, "\r\n"), s, rwc)
}

func test(t *testing.T, snd string, rcv string) (s *smtpd.Session, rwc *MockConn) {
	sndbuf := bytes.NewBuffer([]byte(snd))
	rcvbuf := bytes.NewBuffer(nil)
	rwc = &MockConn{sndbuf, rcvbuf}

	c := textproto.NewConn(rwc)
	assert.NotNil(t, c, "")

	s = smtpd.NewSession(c, nil, nil)
	assert.NotNil(t, s, "")

	s.Serve(make(chan struct{}, 1))

	responses, err := ioutil.ReadAll(rcvbuf)
	assert.NoError(t, err, "")
	assert.Equal(t, rcv, string(responses), "")

	fmt.Println(s)

	return s, rwc
}

func test2(t *testing.T, snd string, rcv string, s *smtpd.Session, rcvbuf io.Reader) *smtpd.Session {
	s.Serve(make(chan struct{}, 1))

	responses, err := ioutil.ReadAll(rcvbuf)
	assert.NoError(t, err, "")
	assert.Equal(t, rcv, string(responses), "")

	return s
}
