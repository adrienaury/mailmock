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
)

// TransactionState is the state of a Transaction.
type TransactionState string

// Transaction States
const (
	TSInitiated  TransactionState = "initiated"    // no command was received yet
	TSInProgress TransactionState = "in progress"  // sender address is set, and recipients are being filled
	TSData       TransactionState = "reading data" // sender and recipients are set, waiting for data to complete
	TSCompleted  TransactionState = "completed"    // data is received, the transaction is successfully completed
	TSAborted    TransactionState = "aborted"      // the transaction is not complete
)

// Transaction represents either a successful, ongoing or aborted SMTP transaction.
type Transaction struct {
	Mail    Mail             `json:"mail"`
	State   TransactionState `json:"state"`
	History []string         `json:"history"`
}

// NewTransaction creates a new SMTP transaction with initial state set to TSInitiated.
func NewTransaction() *Transaction {
	return &Transaction{State: TSInitiated}
}

// Process reads the given command, updates the transaction and returns appropriate response.
func (tr *Transaction) Process(cmd *Command) (*Response, error) {
	if tr != nil {
		tr.History = append(tr.History, cmd.FullCmd)
		r, err := tr.handleCommand(cmd)
		if err != nil {
			tr.History = tr.History[0 : len(tr.History)-1]
		} else {
			tr.History = append(tr.History, r.String())
		}
		return r, err
	}
	return nil, fmt.Errorf("No transaction available to process [%v]", cmd)
}

// Data sets full data, this method can only be user during TSData phase.
func (tr *Transaction) Data(data []string) (*Response, error) {
	if tr != nil && tr.State == TSData {
		tr.History = append(tr.History, data...)
		tr.Mail.Content = data
		tr.State = TSCompleted
		r := r(Success)
		tr.History = append(tr.History, ".")
		tr.History = append(tr.History, r.String())
		return r, nil
	}
	return nil, fmt.Errorf("No transaction available to process data")
}

// Abort sets transaction's state to TSAborted.
func (tr *Transaction) Abort() error {
	if tr != nil && (tr.State == TSInitiated || tr.State == TSInProgress || tr.State == TSData) {
		tr.State = TSAborted
		return nil
	}
	if tr != nil && tr.State == TSCompleted {
		return fmt.Errorf("This transaction is already completed, it can't be aborted")
	}

	if tr != nil && tr.State == TSAborted {
		return nil
	}
	if tr != nil {
		panic("Coding Error : every case must be implemented")
	}
	return nil
}

func (tr *Transaction) handleCommand(cmd *Command) (*Response, error) {
	switch tr.State {
	case TSInitiated:
		return tr.handleCommandInitiated(cmd)
	case TSInProgress:
		return tr.handleCommandInProgress(cmd)
	case TSData:
		return nil, fmt.Errorf("This transaction can't receive commands right now")
	case TSCompleted:
		return tr.handleCommandCompleted(cmd)
	case TSAborted:
		return tr.handleCommandAborted(cmd)
	}
	panic("Coding Error : every case must be implemented")
}

func (tr *Transaction) handleCommandInitiated(cmd *Command) (*Response, error) {
	switch cmd.Name {
	case "MAIL":
		tr.Mail.Envelope.Sender = cmd.NamedArgs["FROM"]
		tr.State = TSInProgress
		return r(Success), nil
	}
	return r(BadSequence), nil
}

func (tr *Transaction) handleCommandInProgress(cmd *Command) (*Response, error) {
	switch cmd.Name {
	case "RCPT":
		tr.Mail.Envelope.Recipients = append(tr.Mail.Envelope.Recipients, cmd.NamedArgs["TO"])
		return r(Success), nil
	case "DATA":
		if len(tr.Mail.Envelope.Recipients) > 0 {
			tr.State = TSData
			return r(Data), nil
		}
	}
	return r(BadSequence), nil
}

func (tr *Transaction) handleCommandCompleted(cmd *Command) (*Response, error) {
	return nil, fmt.Errorf("Sorry, this transaction is completed ans doen't accept any command")
}

func (tr *Transaction) handleCommandAborted(cmd *Command) (*Response, error) {
	return nil, fmt.Errorf("Sorry, this transaction is aborted ans doen't accept any command")
}

func (tr Transaction) String() string {
	return fmt.Sprintf("Transaction %v [%p]", tr.State, &tr)
}
