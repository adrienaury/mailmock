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
	assert.NoError(t, err, "Initiated transactions MUST NOT return an error after a well-formed MAIL command")
	assert.NotNil(t, res, "Initiated transactions MUST return a response to a well-formed MAIL command")
	assert.Equal(t, smtpd.CodeSuccess, res.Code, "Initiated transactions MUST return response code 250 to a well-formed MAIL command")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "Initiated transactions MUST mutate to in progress state after a well-formed MAIL command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "Initiated transactions MUST update their history after a well-formed MAIL command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "Initiated transactions MUST update the sender property after a well-formed MAIL command")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "In progress transactions MUST NOT return an error after a well-formed RCPT command")
	assert.NotNil(t, res, "In progress transactions MUST return a response to a well-formed RCPT command")
	assert.Equal(t, smtpd.CodeSuccess, res.Code, "In progress transactions MUST return response code 250 to a well-formed RCPT command")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "In progress transactions MUST NOT mutate state after a well-formed RCPT command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK"}, tr.History, "In progress transactions MUST update their history after a well-formed RCPT command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "In progress transactions MUST update the recipient property after a well-formed RCPT command")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "In progress transactions MUST NOT return an error after a well-formed DATA command")
	assert.NotNil(t, res, "In progress transactions MUST return a response to a well-formed DATA command")
	assert.Equal(t, smtpd.CodeAskForData, res.Code, "In progress transactions MUST return response code 250 to a well-formed DATA command")
	assert.Equal(t, smtpd.TSData, tr.State, "Initiated transactions MUST mutate to data state after a well-formed DATA command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "In progress transactions MUST update their history after a well-formed DATA command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "In progress transactions MUST NOT update the mail after a well-formed DATA command")

	res, err = tr.Data(MailData)
	assert.NoError(t, err, "Data transactions MUST NOT return an error during data transfer")
	assert.NotNil(t, res, "Data transactions MUST return a response after a successfull data transfer")
	assert.Equal(t, smtpd.CodeSuccess, res.Code, "Data transactions MUST return response code 250 to a well-formed data transfer")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Data transactions MUST mutate to completed state after a well-formed data transfer")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Data transactions MUST update their history after a well-formed data transfer")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Data transactions MUST NOT update the mail content after a well-formed data transfer")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "Completed transactions MUST return an error after any unknown command")
	assert.Nil(t, res, "Completed transactions MUST NOT return a response to any unknown command")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after any unknown command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after any unknown command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after any unknown command")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "Completed transactions MUST return an error after a well-formed MAIL command")
	assert.Nil(t, res, "Completed transactions MUST NOT return a response to a well-formed MAIL command")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after a well-formed MAIL command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after a well-formed MAIL command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after a well-formed MAIL command")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "Completed transactions MUST return an error after a well-formed RCPT command")
	assert.Nil(t, res, "Completed transactions MUST NOT return a response to a well-formed RCPT command")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after a well-formed RCPT command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after a well-formed RCPT command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after a well-formed RCPT command")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "Completed transactions MUST return an error after a well-formed DATA command")
	assert.Nil(t, res, "Completed transactions MUST NOT return a response to a well-formed DATA command")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after a well-formed DATA command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after a well-formed DATA command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after a well-formed DATA command")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "Completed transactions MUST return an error after an attempt to transfer data")
	assert.Nil(t, res, "Completed transactions MUST NOT return a response to an attempt to transfer data")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after an attempt to transfer data")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after an attempt to transfer data")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after an attempt to transfer data")

	err = tr.Abort()
	assert.Error(t, err, "Completed transactions MUST return an error after an attempt to abort transaction")
	assert.Equal(t, smtpd.TSCompleted, tr.State, "Completed transactions MUST NOT mutate state after an attempt to abort transaction")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>", MailData[0], MailData[1], MailData[2], ".", "250 OK"}, tr.History, "Completed transactions MUST NOT update their history after an attempt to abort transaction")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}, Content: MailData}, tr.Mail, "Completed transactions MUST NOT update the mail after an attempt to abort transaction")

	fmt.Println(tr)
}

func TestTransactionAbort1(t *testing.T) {
	tr := smtpd.NewTransaction()

	err := tr.Abort()
	assert.NoError(t, err, "Initiated transactions MUST NOT return an error after an attempt to abort transaction")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Initiated transactions MUST mutate to aborted state after an attempt to abort transaction")
	assert.Empty(t, tr.History, "Initiated transactions MUST NOT update their history after an attempt to abort transaction")
	assert.Empty(t, tr.Mail, "Initiated transactions MUST NOT update the mail after an attempt to abort transaction")

	res, err := tr.Process(&OtherCommand)
	assert.Error(t, err, "Aborted transactions MUST return an error after any unknown command")
	assert.Nil(t, res, "Aborted transactions MUST NOT return a response to any unknown command")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after any unknown command")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after any unknown command")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after any unknown command")

	res, err = tr.Process(&MailCommand)
	assert.Error(t, err, "Aborted transactions MUST return an error after a well-formed MAIL command")
	assert.Nil(t, res, "Aborted transactions MUST NOT return a response to a well-formed MAIL command")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after a well-formed MAIL command")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after a well-formed MAIL command")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after a well-formed MAIL command")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "Aborted transactions MUST return an error after a well-formed RCPT command")
	assert.Nil(t, res, "Aborted transactions MUST NOT return a response to a well-formed RCPT command")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after a well-formed RCPT command")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after a well-formed RCPT command")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after a well-formed RCPT command")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "Aborted transactions MUST return an error after a well-formed DATA command")
	assert.Nil(t, res, "Aborted transactions MUST NOT return a response to a well-formed DATA command")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after a well-formed DATA command")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after a well-formed DATA command")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after a well-formed DATA command")

	res, err = tr.Data(MailData)
	assert.Error(t, err, "Aborted transactions MUST return an error after an attempt to transfer data")
	assert.Nil(t, res, "Aborted transactions MUST NOT return a response to an attempt to transfer data")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after an attempt to transfer data")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after an attempt to transfer data")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after an attempt to transfer data")

	err = tr.Abort()
	assert.NoError(t, err, "Aborted transactions MUST return an error after an attempt to abort transaction")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Aborted transactions MUST NOT mutate state after an attempt to abort transaction")
	assert.Empty(t, tr.History, "Aborted transactions MUST NOT update their history after an attempt to abort transaction")
	assert.Empty(t, tr.Mail, "Aborted transactions MUST NOT update the mail after an attempt to abort transaction")

	fmt.Println(tr)
}

func TestTransactionAbort2(t *testing.T) {
	tr := smtpd.NewTransaction()

	_, err := tr.Process(&MailCommand)
	assert.NoError(t, err)

	err = tr.Abort()
	assert.NoError(t, err, "In progress transactions MUST NOT return an error after an attempt to abort transaction")
	assert.Equal(t, smtpd.TSAborted, tr.State, "In progress transactions MUST mutate to aborted state after an attempt to abort transaction")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK"}, tr.History, "In progress transactions MUST NOT update their history after an attempt to abort transaction")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "In progress transactions MUST NOT update the mail after an attempt to abort transaction")

	fmt.Println(tr)
}

func TestTransactionAbort3(t *testing.T) {
	tr := smtpd.NewTransaction()

	_, err := tr.Process(&MailCommand)
	assert.NoError(t, err)

	_, err = tr.Process(&RcptCommand)
	assert.NoError(t, err)

	_, err = tr.Process(&DataCommand)
	assert.NoError(t, err)

	err = tr.Abort()
	assert.NoError(t, err, "Data transactions MUST NOT return an error after an attempt to abort transaction")
	assert.Equal(t, smtpd.TSAborted, tr.State, "Data transactions MUST mutate to aborted state after an attempt to abort transaction")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "Data transactions MUST NOT update their history after an attempt to abort transaction")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "Data transactions MUST NOT update the mail after an attempt to abort transaction")

	fmt.Println(tr)
}

func TestTransactionWrongSequence1(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&RcptCommand)
	assert.NoError(t, err, "Initiated transactions MUST NOT return an error after a well-formed RCPT command")
	assert.NotNil(t, res, "Initiated transactions MUST return a response to a well-formed RCPT command")
	assert.Equal(t, smtpd.CodeBadSequence, res.Code, "Initiated transactions MUST return response code 503 to a well-formed RCPT command")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "Initiated transactions MUST NOT mutate state after a well-formed RCPT command")
	assert.Equal(t, []string{RcptCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "Initiated transactions MUST update their history after a well-formed RCPT command")
	assert.Empty(t, tr.Mail, "Initiated transactions MUST NOT update the mail after a well-formed RCPT command")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "Initiated transactions MUST NOT return an error after a well-formed DATA command")
	assert.NotNil(t, res, "Initiated transactions MUST return a response to a well-formed DATA command")
	assert.Equal(t, smtpd.CodeBadSequence, res.Code, "Initiated transactions MUST return response code 503 to a well-formed DATA command")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "Initiated transactions MUST NOT mutate state after a well-formed DATA command")
	assert.Equal(t, []string{RcptCommand.FullCmd, "503 Bad sequence of commands", DataCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "Initiated transactions MUST update their history after a well-formed DATA command")
	assert.Empty(t, tr.Mail, "Initiated transactions MUST NOT update the mail after a well-formed DATA command")

	fmt.Println(tr)
}

func TestTransactionWrongSequence2(t *testing.T) {
	tr := smtpd.NewTransaction()

	_, err := tr.Process(&MailCommand)
	assert.NoError(t, err)

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "In progress transactions MUST NOT return an error after a well-formed MAIL command")
	assert.NotNil(t, res, "In progress transactions MUST return a response to a well-formed MAIL command")
	assert.Equal(t, smtpd.CodeBadSequence, res.Code, "In progress transactions MUST return response code 503 to a well-formed MAIL command")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "In progress transactions MUST NOT mutate state after a well-formed MAIL command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", MailCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "In progress transactions MUST update their history after a well-formed MAIL command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "In progress transactions MUST NOT update the mail after a well-formed MAIL command")

	fmt.Println(tr)
}

func TestTransactionMissingRecipient(t *testing.T) {
	tr := smtpd.NewTransaction()

	_, err := tr.Process(&MailCommand)
	assert.NoError(t, err)

	res, err := tr.Process(&DataCommand)
	assert.NoError(t, err, "In progress transaction with no recipients MUST NOT return an error after a well-formed DATA command")
	assert.NotNil(t, res, "In progress transaction with no recipients MUST return a response to a well-formed DATA command")
	assert.Equal(t, smtpd.CodeBadSequence, res.Code, "In progress transaction with no recipients MUST return response code 503 to a well-formed DATA command")
	assert.Equal(t, smtpd.TSInProgress, tr.State, "In progress transaction with no recipients MUST NOT mutate state after a well-formed DATA command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", DataCommand.FullCmd, "503 Bad sequence of commands"}, tr.History, "In progress transactions with no recipients MUST update their history after a well-formed DATA command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com"}}, tr.Mail, "In progress transactions with no recipients MUST NOT update the mail after a well-formed DATA command")

	fmt.Println(tr)
}

func TestTransactionUnexpectedData(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Data(MailData)
	assert.Error(t, err, "Initiated transaction MUST return an error after an attempt to transfer data")
	assert.Nil(t, res, "Initiated transaction MUST NO return a response after an attempt to transfer data")
	assert.Equal(t, smtpd.TSInitiated, tr.State, "Initiated transactions MUST NOT mutate state after an attempt to transfer data")
	assert.Empty(t, tr.History, "Initiated transactions MUST NOT update their history after an attempt to transfer data")
	assert.Empty(t, tr.Mail, "Initiated transactions MUST NOT update the mail after an attempt to transfer data")

	fmt.Println(tr)
}

func TestTransactionUnexpectedCommand(t *testing.T) {
	tr := smtpd.NewTransaction()

	_, err := tr.Process(&MailCommand)
	assert.NoError(t, err)

	_, err = tr.Process(&RcptCommand)
	assert.NoError(t, err)

	_, err = tr.Process(&DataCommand)
	assert.NoError(t, err)

	res, err := tr.Process(&MailCommand)
	assert.Error(t, err, "Data transactions MUST return an error after a well-formed MAIL command")
	assert.Nil(t, res, "Data transactions MUST NOT return a response after a well-formed MAIL command")
	assert.Equal(t, smtpd.TSData, tr.State, "Data transactions MUST NOT mutate state after a well-formed MAIL command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "Data transactions MUST NOT update their history after a well-formed MAIL command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "Data transactions MUST NOT update the mail after a well-formed MAIL command")

	res, err = tr.Process(&RcptCommand)
	assert.Error(t, err, "Data transactions MUST return an error after a well-formed RCPT command")
	assert.Nil(t, res, "Data transactions MUST NOT return a response after a well-formed RCPT command")
	assert.Equal(t, smtpd.TSData, tr.State, "Data transactions MUST NOT mutate state after a well-formed RCPT command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "Data transactions MUST NOT update their history after a well-formed RCPT command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "Data transactions MUST NOT update the mail after a well-formed RCPT command")

	res, err = tr.Process(&DataCommand)
	assert.Error(t, err, "Data transactions MUST return an error after a well-formed DATA command")
	assert.Nil(t, res, "Data transactions MUST NOT return a response after a well-formed DATA command")
	assert.Equal(t, smtpd.TSData, tr.State, "Data transactions MUST NOT mutate state after a well-formed DATA command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "Data transactions MUST NOT update their history after a well-formed DATA command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "Data transactions MUST NOT update the mail after a well-formed DATA command")

	res, err = tr.Process(&OtherCommand)
	assert.Error(t, err, "Data transactions MUST return an error after any unknown command")
	assert.Nil(t, res, "Data transactions MUST NOT return a response after any unknown command")
	assert.Equal(t, smtpd.TSData, tr.State, "Data transactions MUST NOT mutate state after any unknown command")
	assert.Equal(t, []string{MailCommand.FullCmd, "250 OK", RcptCommand.FullCmd, "250 OK", DataCommand.FullCmd, "354 Start mail input; end with <CRLF>.<CRLF>"}, tr.History, "Data transactions MUST NOT update their history after any unknown command")
	assert.Equal(t, smtpd.Mail{Envelope: smtpd.Envelope{Sender: "sender@example.com", Recipients: []string{"recipient@example.com"}}}, tr.Mail, "Data transactions MUST NOT update the mail after any unknown command")

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
	assert.NoError(t, err, "")
	assert.Nil(t, tr, "")

	fmt.Println(tr)
}

func TestTransactionInvalidState(t *testing.T) {
	tr := smtpd.Transaction{State: ""}

	res, err := tr.Data(MailData)
	assert.Error(t, err, "")
	assert.Nil(t, res, "")

	assert.Panics(t, func() { tr.Process(&MailCommand) }, "")

	assert.Panics(t, func() { tr.Abort() }, "")

	fmt.Println(tr)
}
