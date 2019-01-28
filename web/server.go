package web

import (
	"log"
	"net"
	"net/http"
)

// Server represents a local web server.
type Server struct {
	Directory string
	Address   string
}

// ServeStatic starts a static server on address serving contents of directory.
func (s *Server) ServeStatic(verbose bool) (*http.Server, int, error) {
	srv := &http.Server{}

	fs := http.FileServer(http.Dir(s.Directory))
	http.Handle("/", fs)

	port := ":0"
	if s.Address != "" {
		port = s.Address
	}
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, 0, err
	}
	p := listener.Addr().(*net.TCPAddr).Port
	if verbose {
		log.Printf("Serving at http://127.0.0.1:%v\n", p)
	}

	go func() {
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return srv, p, nil
}
