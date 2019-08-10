// Mailmock - Lighweight SMTP server for testing
// Copyright (C) 2019  Adrien Aury
//
// This file is part of Mailmock.
//
// Mailmock is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Mailmock is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with Mailmock.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"fmt"
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

	fmt.Printf("Mailmock v%v  Copyright (C) 2019  Adrien Aury", version)
	fmt.Printf("This program comes with ABSOLUTELY NO WARRANTY; for details type `show w'.")
	fmt.Printf("This is free software, and you are welcome to redistribute it")
	fmt.Printf("under certain conditions; type `show c' for details.")

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
