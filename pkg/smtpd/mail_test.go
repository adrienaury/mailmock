package smtpd_test

import (
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func TestMail(t *testing.T) {
	mail := smtpd.Mail{}
	assert.Equal(t, "MAIL FROM:\nRCPT TO:\n", mail.String(), "Invalid mail string representation")
}
