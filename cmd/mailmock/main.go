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
	"github.com/adrienaury/mailmock/internal/log"
	"github.com/adrienaury/mailmock/internal/repository"
	"github.com/adrienaury/mailmock/pkg/smtpd"
	"github.com/heptio/workgroup"
	"github.com/sirupsen/logrus"
	logur "logur.dev/adapter/logrus"
)

// Provisioned by ldflags
// nolint: gochecknoglobals
var (
	version   string
	commit    string
	buildDate string
	builtBy   string
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

	// sets the SMTP greeting banner
	smtpd.SetReply(smtpd.CodeReady,
		fmt.Sprintf("<domain> Mailmock %v Service ready - this is a testing SMTP server, it does not deliver e-mails", version))

	// logrus initialization
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)

	logger := log.NewLoggerAdapter(logur.New(logrus.StandardLogger()))
	logger = logger.WithFields(log.Fields{
		log.FieldApp: "mailmock",
	})
	logger.Debug("Build information", log.Fields{
		log.FieldVersion:   version,
		log.FieldCommit:    commit,
		log.FieldBuildDate: buildDate,
		log.FieldBuiltBy:   builtBy,
	})

	loggerSMTP := logger.WithFields(log.Fields{
		log.FieldService: "smtp",
	})

	loggerHTTP := logger.WithFields(log.Fields{
		log.FieldService: "http",
	})

	group := &workgroup.Group{}
	group.Add(func(stop <-chan struct{}) error {
		// interrupt/kill signals sent from terminal or host on shutdown
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-stop:
			return fmt.Errorf("shutting down OS signal watcher on workgroup stop")
		case i := <-interrupt:
			logger.Info(fmt.Sprintf("Received OS signal %s; beginning shutdown...", i))
			return nil
		}
	})
	group.Add(func(stop <-chan struct{}) error {
		smtpsrv := smtpd.NewServer("main", listenAddr, smtpPort, &th, loggerSMTP)
		return smtpsrv.ListenAndServe(stop)
	})
	group.Add(func(stop <-chan struct{}) error {
		httpsrv := httpd.NewServer("main", listenAddr, httpPort, loggerHTTP)
		return httpsrv.ListenAndServe(stop)
	})
	err := group.Run()
	if err != nil {
		logger.Error("Program exited with error", log.Fields{log.FieldError: err})
		os.Exit(1)
	}
}
