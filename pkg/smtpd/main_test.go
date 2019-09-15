package smtpd_test

import (
	"fmt"
	"net/smtp"
	"os"
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

var th smtpd.TransactionHandler = func(tr *smtpd.Transaction) {
	fmt.Println(tr)
}

func TestMain(m *testing.M) {
	smtpd.SetReply(smtpd.Ready, "Service ready")
	smtpd.SetReply(smtpd.Closing, "Service closing transmission channel")
	srv := smtpd.NewServer("mockmail", "localhost", "1024", &th, nil)
	go func() {
		if err := srv.ListenAndServe(make(chan struct{})); err != nil {
			panic(err)
		}
	}()
	os.Exit(m.Run())
}

func TestNominal(t *testing.T) {
	c, err := smtp.Dial("127.0.0.1:1024")
	assert.NoError(t, err, "Can't contact SMTP server")
	assert.NotNil(t, c, "No connection to SMTP server")

	err = c.Mail("sender@example.org")
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")

	err = c.Rcpt("recipient@example.net")
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")

	wc, err := c.Data()
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")
	assert.NotNil(t, wc, "No data writer")

	_, err = fmt.Fprintf(wc, "Subject: test\n\nThis is the email body")
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")

	err = wc.Close()
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")

	err = c.Quit()
	assert.NoError(t, err, "SMTP server MUST NOT return an error to a valid transaction")
}
