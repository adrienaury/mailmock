package smtpd_test

import (
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func testOk(t *testing.T, input string, exname string, expargs []string, exnargs map[string]string) {
	cmd, res := smtpd.ParseCommand(input)
	assert.NotNil(t, cmd)
	assert.Equal(t, input, cmd.FullCmd)
	assert.Equal(t, exname, cmd.Name)
	assert.Equal(t, expargs, cmd.PositionalArgs)
	assert.Equal(t, exnargs, cmd.NamedArgs)
	assert.Nil(t, res)
}

func testKo(t *testing.T, input string, code int16, msg string) {
	cmd, res := smtpd.ParseCommand(input)
	assert.Nil(t, cmd)
	assert.NotNil(t, res)
	assert.Equal(t, code, res.Code)
	assert.Equal(t, msg, res.Msg)
}

func TestCommandNominal(t *testing.T) {
	testOk(t, "HELO localhost", "HELO", []string{"localhost"}, map[string]string{})
	testOk(t, "EHLO localhost", "EHLO", []string{"localhost"}, map[string]string{})
	testOk(t, "MAIL FROM:<sender@example.com>", "MAIL", []string{}, map[string]string{"FROM": "<sender@example.com>"})
	testOk(t, "RCPT TO:<recipient@example.com>", "RCPT", []string{}, map[string]string{"TO": "<recipient@example.com>"})
	testOk(t, "DATA", "DATA", []string{}, map[string]string{})
	testOk(t, "NOOP", "NOOP", []string{}, map[string]string{})
	testOk(t, "RSET", "RSET", []string{}, map[string]string{})
	testOk(t, "QUIT", "QUIT", []string{}, map[string]string{})
	testOk(t, "VRFY", "VRFY", []string{}, map[string]string{})
}

func TestCommandLowercase(t *testing.T) {
	testOk(t, "helo localhost", "HELO", []string{"localhost"}, map[string]string{})
	testOk(t, "ehlo localhost", "EHLO", []string{"localhost"}, map[string]string{})
	testOk(t, "mail from:<sender@example.com>", "MAIL", []string{}, map[string]string{"FROM": "<sender@example.com>"})
	testOk(t, "rcpt to:<recipient@example.com>", "RCPT", []string{}, map[string]string{"TO": "<recipient@example.com>"})
	testOk(t, "data", "DATA", []string{}, map[string]string{})
	testOk(t, "noop", "NOOP", []string{}, map[string]string{})
	testOk(t, "rset", "RSET", []string{}, map[string]string{})
	testOk(t, "quit", "QUIT", []string{}, map[string]string{})
	testOk(t, "vrfy", "VRFY", []string{}, map[string]string{})
}

func TestCommandNumberArguments1(t *testing.T) {
	testKo(t, "HELO localhost test", 501, "Syntax error in parameters or arguments")
	testKo(t, "EHLO localhost test", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL FROM:<sender@example.com> test", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT TO:<recipient@example.com> test", 501, "Syntax error in parameters or arguments")
	testKo(t, "DATA test", 501, "Syntax error in parameters or arguments")
	testOk(t, "NOOP test", "NOOP", []string{}, map[string]string{})
	testKo(t, "RSET test", 501, "Syntax error in parameters or arguments")
	testKo(t, "QUIT test", 501, "Syntax error in parameters or arguments")
	testKo(t, "VRFY test", 501, "Syntax error in parameters or arguments")
}

func TestCommandNumberArguments2(t *testing.T) {
	testKo(t, "HELO localhost TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "EHLO localhost TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL FROM:<sender@example.com> TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT TO:<recipient@example.com> TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "DATA TEST:test", 501, "Syntax error in parameters or arguments")
	testOk(t, "NOOP TEST:test", "NOOP", []string{}, map[string]string{})
	testKo(t, "RSET TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "QUIT TEST:test", 501, "Syntax error in parameters or arguments")
	testKo(t, "VRFY TEST:test", 501, "Syntax error in parameters or arguments")
}

func TestCommandNumberArguments3(t *testing.T) {
	testKo(t, "HELO", 501, "Syntax error in parameters or arguments")
	testKo(t, "EHLO", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL test", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT test", 501, "Syntax error in parameters or arguments")
}

func TestCommandWrongArgument(t *testing.T) {
	testKo(t, "MAIL MORF:<sender@example.com>", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL <sender@example.com>", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT OT:<recipient@example.com>", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT <recipient@example.com>", 501, "Syntax error in parameters or arguments")
}

func TestCommandEmptyArgument(t *testing.T) {
	testKo(t, "HELO ", 501, "Syntax error in parameters or arguments")
	testKo(t, "EHLO ", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL FROM:", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT TO:", 501, "Syntax error in parameters or arguments")
}

func TestCommandWrongName(t *testing.T) {
	testKo(t, "FAKE", 500, "Syntax error, command unrecognized")
	testKo(t, "FAKE test", 500, "Syntax error, command unrecognized")
	testKo(t, "FAKE TEST:test", 500, "Syntax error, command unrecognized")
}
