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

// Package httpd exposes the REST API of Mailmock.
package httpd

import (
	"context"
	"net"
	"net/http"

	"github.com/adrienaury/mailmock/internal/log"
)

// Server is holding the HTTP server properties.
type Server struct {
	name   string
	host   string
	port   string
	logger log.Logger
}

// NewServer creates a HTTP server.
func NewServer(name string, host string, port string, logger log.Logger) *Server {
	if logger == nil {
		logger = log.DefaultLogger
	}
	l := logger.WithFields(log.Fields{
		log.FieldServer: name,
		log.FieldListen: net.JoinHostPort(host, port),
	})
	return &Server{name, host, port, l}
}

// ListenAndServe starts listening for clients connection and serves requests.
func (srv *Server) ListenAndServe(stop <-chan struct{}) error {
	router := srv.Routes()

	s := http.Server{
		Addr:    net.JoinHostPort(srv.host, srv.port),
		Handler: router,
	}

	go func() {
		<-stop // wait for stop signal
		s.Shutdown(context.Background())
	}()

	srv.logger.Info("HTTP Server is listening")
	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		srv.logger.Error("HTTP Server failed", log.Fields{log.FieldError: err})
		return err
	}
	srv.logger.Info("HTTP Server is stopped")
	return nil
}
