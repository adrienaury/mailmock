package main

import (
	"os"

	"github.com/adrienaury/mailmock/internal/httpd"
	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/adrienaury/mailmock/pkg/smtpd"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version string
	commit  string
	build   string
	builtBy string
)

var th smtpd.TransactionHandler = func(tr *smtpd.Transaction) {
	repository.Store(tr)
}

func main() {

	defaultSMTPPort := "smtp"
	defaultHTTPPort := "http"
	defaultListenAddr := ""

	smtpPort := os.Getenv("MAILMOCK_SMTP_PORT")
	if smtpPort == "" {
		smtpPort = defaultSMTPPort
	}

	httpPort := os.Getenv("MAILMOCK_HTTP_PORT")
	if httpPort == "" {
		httpPort = defaultHTTPPort
	}

	listenAddr := os.Getenv("MAILMOCK_LISTEN_ADDR")
	if listenAddr == "" {
		listenAddr = defaultListenAddr
	}

	// starts the SMTP server
	smtpsrv := smtpd.NewServer("mailmock", listenAddr, smtpPort, &th)
	go smtpsrv.ListenAndServe()

	httpsrv := httpd.NewServer("mailmock", listenAddr, httpPort)
	httpsrv.ListenAndServe()

}
