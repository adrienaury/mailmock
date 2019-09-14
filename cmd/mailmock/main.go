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
	"flag"
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
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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

	var cfgFile string

	flag.String("httpPort", "http", "HTTP Port")
	flag.String("smtpPort", "smtp", "SMTP Port")
	flag.String("address", "", "Listening address")
	flag.String("logLevel", "info", "Log level (trace, debug, info, warn, error)")
	flag.StringVar(&cfgFile, "config", "", "Configuration file")

	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	viper.SetEnvPrefix("mailmock") // will be uppercased automatically

	viper.BindEnv("httpPort")
	viper.BindEnv("smtpPort")
	viper.BindEnv("address")
	viper.BindEnv("logLevel")

	viper.SetDefault("httpPort", "http")
	viper.SetDefault("smtpPort", "smtp")
	viper.SetDefault("address", "")
	viper.SetDefault("logLevel", "info")

	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.SetConfigName("config")          // name of config file (without extension)
		viper.AddConfigPath("/etc/mailmock/")  // path to look for the config file in
		viper.AddConfigPath("$HOME/.mailmock") // call multiple times to add many search paths
		viper.AddConfigPath(".")               // optionally look for config in the working directory
	}

	viper.BindPFlags(pflag.CommandLine)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("failed to read config file: %s", err))
		}
	}

	smtpPort := viper.GetString("smtpPort")
	httpPort := viper.GetString("httpPort")
	listenAddr := viper.GetString("address")
	logLevel, err := logrus.ParseLevel(viper.GetString("logLevel"))
	if err != nil {
		panic(err)
	}

	// sets the SMTP greeting banner
	smtpd.SetReply(smtpd.Ready,
		fmt.Sprintf("<domain> Mailmock %v Service ready", version),
		"This is a testing SMTP server, it does not deliver e-mails")

	// logrus initialization
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logLevel)

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
	err = group.Run()
	if err != nil {
		logger.Error("Program exited with error", log.Fields{log.FieldError: err})
		os.Exit(1)
	}
}
