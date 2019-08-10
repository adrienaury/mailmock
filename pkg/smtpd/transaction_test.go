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
	"fmt"
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

var (
	MailCommand  smtpd.Command = smtpd.Command{Name: "MAIL", NamedArgs: map[string]string{"FROM": "sender@example.com"}, FullCmd: "MAIL FROM:<sender@example.com>"}
	RcptCommand  smtpd.Command = smtpd.Command{Name: "RCPT", NamedArgs: map[string]string{"TO": "recipient@example.com"}, FullCmd: "RCPT TO:<recipient@example.com>"}
	DataCommand  smtpd.Command = smtpd.Command{Name: "DATA", FullCmd: "DATA"}
	MailData     []string      = []string{"Subject: test", "", "This is the email body"}
	OtherCommand smtpd.Command = smtpd.Command{Name: "FAKE", FullCmd: "FAKE"}
)

func TestTransactionNominal(t *testing.T) {
	tr := smtpd.NewTransaction()
	assert.Equal(t, smtpd.TSInitiated, tr.State, "A newly created transaction MUST have an initiated State")
	assert.Empty(t, tr.History, "A newly created transaction MUST have an empty History")
	assert.Empty(t, tr.Mail, "A newly created transaction MUST have an empty Mail")

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "A transaction in initiated state MUST NOT return an error to a well-formed MAIL command")
	assert.NotNil(t, res, "A transaction in initiated state MUST return a response to a well-formed MAIL command")
	assert.Equal(t, int16(250), res.Code, "A transaction in initiated state MUST return with response code 250 to a well-formed MAIL command")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "A transaction in initiated state MUST change to in progress state after a well-formed MAIL command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "A transaction in initiated state MUST update its history after a well-formed MAIL command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "A transaction in initiated state MUST update the sender after a well-formed MAIL command")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(354), res.Code, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionAbort1(t *testing.T) {
	tr := smtpd.NewTransaction()
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	err := tr.Abort()
	assert.NoError(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err := tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionAbort2(t *testing.T) {
	tr := smtpd.NewTransaction()
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	err = tr.Abort()
	assert.NoError(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionAbort3(t *testing.T) {
	tr := smtpd.NewTransaction()
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	err = tr.Abort()
	assert.NoError(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionAbort4(t *testing.T) {
	tr := smtpd.NewTransaction()
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(354), res.Code, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	err = tr.Abort()
	assert.NoError(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Equal(t, smtpd.TSAborted, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionWrongSequence1(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Equal(t, []string{RcptCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Equal(t, []string{RcptCommand.FullCmd, "503 Bad sequence of commands", DataCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionWrongSequence2(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", MailCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionMissingRecipient(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", DataCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionUnexpectedData(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "")
	assert.Empty(t, tr.History, "")
	assert.Empty(t, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionUnexpectedCommand(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(354), res.Code, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Equal(t, smtpd.TSData, tr.State, "")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "")

	fmt.Println(tr)
}

func TestTransactionNil(t *testing.T) {
	var tr *smtpd.Transaction

	res, err := tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Nil(t, tr, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")
	assert.Nil(t, tr, "")

	err = tr.Abort()
	assert.Error(t, err, "")
	assert.Nil(t, tr, "")

	fmt.Println(tr)
}

func TestTransactionInvalidState(t *testing.T) {
	tr := smtpd.Transaction{State: ""}

	res, err := tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")

	err = tr.Abort()
	assert.Error(t, err, "")

	fmt.Println(tr)
}
