package smtpd_test

import (
	"fmt"
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

var (
	MailCommand smtpd.Command = smtpd.Command{Name: "MAIL", NamedArgs: map[string]string{"FROM": "sender@example.com"}}
	RcptCommand smtpd.Command = smtpd.Command{Name: "RCPT", NamedArgs: map[string]string{"TO": "recipient@example.com"}}
	DataCommand smtpd.Command = smtpd.Command{Name: "DATA"}
	MailData    string        = "Subject: test\n\nThis is the email body"
)

func TestNominal(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")

	res, err = tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(354), res.Code, "")

	res, err = tr.Data(MailData)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")

	fmt.Println(tr)
}

func TestWrongSequence1(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&RcptCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")

	fmt.Println(tr)
}

func TestWrongSequence2(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")

	fmt.Println(tr)
}

func TestMissingRecipient(t *testing.T) {
	tr := smtpd.NewTransaction()

	res, err := tr.Process(&MailCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(250), res.Code, "")

	res, err = tr.Process(&DataCommand)
	assert.NoError(t, err, "")
	assert.NotNil(t, res, "")
	assert.Equal(t, int16(503), res.Code, "")

	fmt.Println(tr)
}
