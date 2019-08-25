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
//
// Linking this library statically or dynamically with other modules is
// making a combined work based on this library.  Thus, the terms and
// conditions of the GNU General Public License cover the whole
// combination.
//
// As a special exception, the copyright holders of this library give you
// permission to link this library with independent modules to produce an
// executable, regardless of the license terms of these independent
// modules, and to copy and distribute the resulting executable under
// terms of your choice, provided that you also meet, for each linked
// independent module, the terms and conditions of the license of that
// module.  An independent module is a module which is not derived from
// or based on this library.  If you modify this library, you may extend
// this exception to your version of the library, but you are not
// obligated to do so.  If you do not wish to do so, delete this
// exception statement from your version.

// Package smtpd contains source code of the SMTP server of Mailmock.
package smtpd

import (
	"net"
	"net/textproto"
	"sync"
	"time"

	"github.com/adrienaury/mailmock/pkg/smtpd/log"
	"github.com/goph/logur"
)

// Server is holding the SMTP server properties.
type Server struct {
	name      string
	host      string
	port      string
	th        *TransactionHandler
	logger    log.Logger
	waitGroup *sync.WaitGroup
}

// NewServer creates a SMTP server.
func NewServer(name string, host string, port string, th *TransactionHandler, logger log.Logger) *Server {
	if logger == nil {
		logger = log.DefaultLogger
	}
	l := logur.WithFields(logger, log.Fields{
		log.FieldServer: name,
		log.FieldListen: net.JoinHostPort(host, port),
	})
	srv := &Server{name, host, port, th, l, &sync.WaitGroup{}}
	return srv
}

// ListenAndServe starts listening for clients connection and serves SMTP commands.
func (srv *Server) ListenAndServe(stop <-chan struct{}) error {
	defer srv.un(srv.trace("ListenAndServe"))
	laddr, err := net.ResolveTCPAddr("tcp", net.JoinHostPort(srv.host, srv.port))
	if nil != err {
		srv.logger.Error("SMTP Server failed to start", log.Fields{log.FieldError: err})
		return err
	}
	ln, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		srv.logger.Error("SMTP Server failed to start", log.Fields{log.FieldError: err})
		return err
	}
	srv.logger.Info("SMTP Server is listening")
	srv.waitGroup.Add(1)
	srv.serve(ln, stop)
	srv.waitGroup.Wait()
	srv.logger.Info("SMTP Server is stopped")
	return nil
}

func (srv *Server) serve(ln *net.TCPListener, stop <-chan struct{}) {
	defer srv.un(srv.trace("serve"))
	defer ln.Close()
	defer srv.waitGroup.Done()
	for {
		select {
		case <-stop:
			return
		default:
		}
		ln.SetDeadline(time.Now().Add(1e9)) // 1 second
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Timeout() {
				continue
			}
			srv.logger.Error("SMTP Server failed to accept connection", log.Fields{log.FieldError: err})
			continue
		}
		srv.waitGroup.Add(1)
		go srv.handleConnection(conn.(*net.TCPConn), stop)
	}
}

func (srv *Server) handleConnection(conn *net.TCPConn, stop <-chan struct{}) {
	defer srv.un(srv.trace("handleConnection"))
	tpc := textproto.NewConn(conn)
	defer tpc.Close()
	defer srv.waitGroup.Done()

	s := NewSession(tpc, srv.th, srv.logger)
	s.tcpConn = conn
	s.Serve(stop)
}

func (srv *Server) trace(s string) string {
	srv.logger.Trace("Entering function", log.Fields{"func": s})
	return s
}

func (srv *Server) un(s string) {
	srv.logger.Trace("Leaving function", log.Fields{"func": s})
}
