package main

import (
	"github.com/adrienaury/mailmock/internal/httpd"
	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/adrienaury/mailmock/pkg/smtpd"
)

var th smtpd.TransactionHandler = func(tr *smtpd.Transaction) {
	repository.Store(tr)
}

func main() {

	// starts the SMTP server
	smtpsrv := smtpd.NewServer("mailmock", "localhost", "smtp", &th)
	go smtpsrv.ListenAndServe()

	httpsrv := httpd.NewServer("mailmock", "localhost", "http")
	httpsrv.ListenAndServe()

}
