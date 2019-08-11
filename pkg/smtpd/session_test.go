package smtpd_test

import (
	"bytes"
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

func test(t *testing.T, snd string, rcv string) {
	sndbuf := bytes.NewBuffer([]byte(snd))
	rcvbuf := bytes.NewBuffer(nil)
	rwc := &MockConn{sndbuf, rcvbuf}

	c := textproto.NewConn(rwc)
	assert.NotNil(t, c, "")

	s := smtpd.NewSession(c, nil)
	assert.NotNil(t, s, "")

	s.Serve()

	responses, err := ioutil.ReadAll(rcvbuf)
	assert.NoError(t, err, "")
	assert.Equal(t, rcv, string(responses), "")
}
