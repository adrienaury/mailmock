package smtpd_test

import (
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	response := smtpd.Response{250, "OK"}
	assert.Equal(t, "250 OK", response.String(), "")
	assert.Equal(t, false, response.IsError(), "")
	assert.Equal(t, true, response.IsSuccess(), "")

	response = smtpd.Response{500, "Not OK"}
	assert.Equal(t, "500 Not OK", response.String(), "")
	assert.Equal(t, true, response.IsError(), "")
	assert.Equal(t, false, response.IsSuccess(), "")
}
