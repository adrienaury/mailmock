package httpd

import (
	"log"
	"net"
	"net/http"

	"github.com/go-chi/chi"
)

// Server is holding the HTTP server properties
type Server struct {
	name string
	host string
	port string
}

// NewServer creates a HTTP server
func NewServer(name string, host string, port string) *Server {
	return &Server{name, host, port}
}

// ListenAndServe starts listening for clients connection and serves requests
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
