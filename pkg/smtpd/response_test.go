package smtpd_test

import (
	"testing"

	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/stretchr/testify/assert"
)

func TestResponse(t *testing.T) {
	response := smtpd.Response{Code: smtpd.Code(250), Msg: "OK"}
	assert.Equal(t, "250 OK", response.String(), "Response is not parsed correctly")
	assert.Equal(t, false, response.IsError(), "Response code 250 indicates a success")
	assert.Equal(t, true, response.IsSuccess(), "Response code 250 indicates a success")

	response = smtpd.Response{Code: smtpd.Code(500), Msg: "Not OK"}
	assert.Equal(t, "500 Not OK", response.String(), "Response is not parsed correctly")
	assert.Equal(t, true, response.IsError(), "Response code 500 indicates a failure")
	assert.Equal(t, false, response.IsSuccess(), "Response code 500 indicates a failure")
}
