package smtpd

import (
	"fmt"
	"strings"
)

// Response holds a 3 digit code and a messsage
type Response struct {
	Code int16
	Msg  string
}

// IsError returns true if the response is an error
func (e Response) IsError() bool {
	return strings.HasPrefix(e.String(), "5")
}

// IsSuccess returns true if the response is an success
func (e Response) IsSuccess() bool {
	return !e.IsError()
}

func (e Response) String() string {
	return fmt.Sprintf("%3d %s", e.Code, e.Msg)
}
