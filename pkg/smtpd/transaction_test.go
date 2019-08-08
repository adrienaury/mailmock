package smtpd_test

import (
	"fmt"
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func TestNominal(t *testing.T) {
	tr := smtpd.NewTransaction()

	command, res := smtpd.ParseCommand("MAIL FROM:<sender@example.com>")
	assert.Nil(t, res, "")
	res, err := tr.Update(command)
	assert.NoError(t, err, "")

	command, res = smtpd.ParseCommand("RCPT TO:<recipient@example.com>")
	assert.Nil(t, res, "")
	res, err = tr.Update(command)
	assert.NoError(t, err, "")

	command, res = smtpd.ParseCommand("DATA")
	assert.Nil(t, res, "")
	res, err = tr.Update(command)
	assert.NoError(t, err, "")

	res, err = tr.Data("Subject: test\n\nThis is the email body")
	assert.NoError(t, err, "")

	fmt.Println(tr)
}
