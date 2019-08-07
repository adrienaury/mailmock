package main

import (
	"github.com/adrienaury/mailmock/internal/httpd"
	"github.com/adrienaury/mailmock/pkg/smtpd"
)

func main() {

	// starts the SMTP server
	smtpsrv := smtpd.NewServer("mailmock", "localhost", "smtp")
	go smtpsrv.ListenAndServe()

	httpsrv := httpd.NewServer("mailmock", "localhost", "http")
	httpsrv.ListenAndServe()

}
