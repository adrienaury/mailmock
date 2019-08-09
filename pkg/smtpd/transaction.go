package smtpd

import (
	"fmt"
	"strings"
)

// TransactionState is the state of a Transaction
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

// NewTransaction creates a new SMTP transaction with initial state set to TSInitiated
func NewTransaction() *Transaction {
	return &Transaction{State: TSInitiated}
}

// Process reads the given command, updates the transaction and returns appropriate response
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
	return nil, fmt.Errorf("TODO")
}

// Data sets full data, this method can only be user during TSData phase
func (tr *Transaction) Data(data string) (*Response, error) {
	if tr != nil && tr.State == TSData {
		tr.History = append(tr.History, data)
		tr.Mail.Content = &data
		tr.State = TSCompleted
		r := &Response{250, "OK"}
		tr.History = append(tr.History, ".")
		tr.History = append(tr.History, r.String())
		return r, nil
	}
	return nil, fmt.Errorf("TODO")
}

// Abort sets transaction's state to TSAborted
func (tr *Transaction) Abort() error {
	if tr != nil && (tr.State == TSInitiated || tr.State == TSInProgress || tr.State == TSData) {
		tr.State = TSAborted
		return nil
	}
	return fmt.Errorf("TODO")
}

func (tr *Transaction) handleCommand(cmd *Command) (*Response, error) {
	switch tr.State {
	case TSInitiated:
		return tr.handleCommandInitiated(cmd)
	case TSInProgress:
		return tr.handleCommandInProgress(cmd)
	case TSData:
		return nil, fmt.Errorf("TODO")
	case TSCompleted:
		return tr.handleCommandCompleted(cmd)
	case TSAborted:
		return tr.handleCommandAborted(cmd)
	}
	return nil, fmt.Errorf("TODO")
}

func (tr *Transaction) handleCommandInitiated(cmd *Command) (*Response, error) {
	switch cmd.Name {
	case "MAIL":
		tr.Mail.Envelope.Sender = cmd.NamedArgs["FROM"]
		tr.State = TSInProgress
		return &Response{250, "OK"}, nil
	}
	return &Response{503, "Bad sequence of commands"}, nil
}

func (tr *Transaction) handleCommandInProgress(cmd *Command) (*Response, error) {
	switch cmd.Name {
	case "RCPT":
		tr.Mail.Envelope.Recipients = append(tr.Mail.Envelope.Recipients, cmd.NamedArgs["TO"])
		return &Response{250, "OK"}, nil
	case "DATA":
		if len(tr.Mail.Envelope.Recipients) > 0 {
			tr.State = TSData
			return &Response{354, "Start mail input; end with <CRLF>.<CRLF>"}, nil
		}
	}
	return &Response{503, "Bad sequence of commands"}, nil
}

func (tr *Transaction) handleCommandCompleted(cmd *Command) (*Response, error) {
	return nil, fmt.Errorf("TODO")
}

func (tr *Transaction) handleCommandAborted(cmd *Command) (*Response, error) {
	return nil, fmt.Errorf("TODO")
}

func (tr Transaction) String() string {
	return fmt.Sprintf("Transaction %v:\n%v", tr.State, strings.Join(tr.History, "\n"))
}
