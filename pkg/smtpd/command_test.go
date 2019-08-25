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
package smtpd_test

import (
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func testOk(t *testing.T, input string, exname string, expargs []string, exnargs map[string]string) {
	cmd, res := smtpd.ParseCommand(input)
	assert.NotNil(t, cmd, "Parsing a valid command MUST NOT return a nil value")
	assert.Equal(t, input, cmd.FullCmd, "FullCmd MUST contain the original full command")
	assert.Equal(t, exname, cmd.Name, "Parsed command name is invalid")
	assert.Equal(t, expargs, cmd.PositionalArgs, "Parsed positional arguments are invalid")
	assert.Equal(t, exnargs, cmd.NamedArgs, "Parsed named arguments are invalid")
	assert.Nil(t, res, "Parsing a valid command MUST return a nil response")
}

func testKo(t *testing.T, input string, code smtpd.Code, msg string) {
	cmd, res := smtpd.ParseCommand(input)
	assert.Nil(t, cmd, "Parsing an invalid command MUST return nil")
	assert.NotNil(t, res, "Parsing an invalid command MUST return a non nil error response")
	assert.Equal(t, code, res.Code, "Response code is not valid")
	assert.Equal(t, msg, res.Msg, "Response message is not valid")
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
	testOk(t, "VRFY test", "VRFY", []string{"test"}, map[string]string{})
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
	testOk(t, "vrfy test", "VRFY", []string{"test"}, map[string]string{})
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
	testKo(t, "VRFY test test", 501, "Syntax error in parameters or arguments")
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
	testKo(t, "VRFY test TEST:test", 501, "Syntax error in parameters or arguments")
}

func TestCommandNumberArguments3(t *testing.T) {
	testKo(t, "HELO", 501, "Syntax error in parameters or arguments")
	testKo(t, "EHLO", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL test", 501, "Syntax error in parameters or arguments")
	testKo(t, "MAIL", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT", 501, "Syntax error in parameters or arguments")
	testKo(t, "RCPT test", 501, "Syntax error in parameters or arguments")
	testKo(t, "VRFY", 501, "Syntax error in parameters or arguments")
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
