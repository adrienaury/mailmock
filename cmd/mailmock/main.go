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
	"os/signal"
	"syscall"

	"github.com/adrienaury/mailmock/internal/httpd"
	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/goph/logur"
	"github.com/goph/logur/adapters/logrusadapter"
	"github.com/sirupsen/logrus"
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

	fmt.Printf(`
     __  __       _ _                      _
    |  \/  |     (_) |                    | |
    | \  / | __ _ _| |_ __ ___   ___   ___| | __
    | |\/| |/ _' | | | '_ ' _ \ / _ \ / __| |/ /
    | |  | | (_| | | | | | | | | (_) | (__|   <
    |_|  |_|\__,_|_|_|_| |_| |_|\___/ \___|_|\_\
        `)

	fmt.Printf("%v - Copyright (C) 2019  Adrien Aury\n\n", version)
	fmt.Println("This program is licensed under the terms of the GNU General Public License v3 (https://www.gnu.org/licenses/gpl-3.0.html)")
	fmt.Println("Source code and documentation are available at https://github.com/adrienaury/mailmock")
	fmt.Println()

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

	hostname := "localhost"
	hostname, _ = os.Hostname()

	// sets the SMTP greeting banner
	smtpd.SetReply(smtpd.CodeReady,
		fmt.Sprintf("%v Mailmock %v Service ready - this is a testing SMTP server, it does not deliver e-mails", hostname, version))

	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logger := logur.WithFields(logrusadapter.New(logrus.StandardLogger()), logur.Fields{
		"app": "mailmock",
	})

	loggerSMTP := logur.WithFields(logger, logur.Fields{
		"service": "smtp",
	})

	loggerHTTP := logur.WithFields(logger, logur.Fields{
		"service": "http",
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan error, 2)
	stop := make(chan struct{})

	go func() {
		<-sigs
		stop <- struct{}{}
	}()

	go func() {
		smtpsrv := smtpd.NewServer("main", listenAddr, smtpPort, &th, loggerSMTP)
		done <- smtpsrv.ListenAndServe(stop)
	}()

	go func() {
		httpsrv := httpd.NewServer("main", listenAddr, httpPort, loggerHTTP)
		done <- httpsrv.ListenAndServe(stop)
	}()

	var stopped, errored bool
	for i := 0; i < cap(done); i++ {
		if err := <-done; err != nil {
			errored = true
		}
		if !stopped {
			stopped = true
			close(stop)
		}
	}
	if errored {
		os.Exit(1)
	}
}
