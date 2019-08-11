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
	srv := smtpd.NewServer("mockmail", "localhost", "smtp", nil)
	go srv.ListenAndServe()
	os.Exit(m.Run())
}

func TestNominal(t *testing.T) {
	c, err := smtp.Dial("127.0.0.1:25")
	assert.NoError(t, err, "")
	assert.NotNil(t, c, "")

	err = c.Mail("sender@example.org")
	assert.NoError(t, err, "")

	err = c.Rcpt("recipient@example.net")
	assert.NoError(t, err, "")

	wc, err := c.Data()
	assert.NoError(t, err, "")
	assert.NotNil(t, wc, "")

	_, err = fmt.Fprintf(wc, "Subject: test\n\nThis is the email body")
	assert.NoError(t, err, "")

	err = wc.Close()
	assert.NoError(t, err, "")

	err = c.Quit()
	assert.NoError(t, err, "")
}
