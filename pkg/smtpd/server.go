package smtpd

import (
	"log"
	"net"
	"net/textproto"
)

// Server is holding the SMTP server properties
type Server struct {
	name string
	host string
	port string
}

// NewServer creates a SMTP server
func NewServer(name string, host string, port string) *Server {
	return &Server{name, host, port}
}

// ListenAndServe starts listening for clients connection and serves SMTP commands
func (srv *Server) ListenAndServe() {
	ln, err := net.Listen("tcp", net.JoinHostPort(srv.host, srv.port))
	if err != nil {
		panic(err)
	}
	srv.serve(ln)
}

func (srv *Server) serve(ln net.Listener) {
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go srv.handleConnection(conn)
	}
}

func (srv *Server) handleConnection(conn net.Conn) {
	tpc := textproto.NewConn(conn)
	defer tpc.Close()

	//s := NewSession(tpc, &srv.mhdl, &srv.thdl)
	//s.Serve()
}
