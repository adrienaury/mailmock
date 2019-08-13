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

// Package httpd exposes the REST API of Mailmock
package httpd

import (
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
)

// Server is holding the HTTP server properties.
type Server struct {
	name string
	host string
	port string
}

// NewServer creates a HTTP server.
func NewServer(name string, host string, port string) *Server {
	return &Server{name, host, port}
}

// ListenAndServe starts listening for clients connection and serves requests.
func (srv *Server) ListenAndServe() {
	router := Routes()

	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}

	if err := chi.Walk(router, walkFunc); err != nil {
		log.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
	}

	log.Fatal(http.ListenAndServe(net.JoinHostPort(srv.host, srv.port), router))
}
