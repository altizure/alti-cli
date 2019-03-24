package web

import (
	"log"
	"net"
	"net/http"
	"strings"
)

// Server represents a local web server.
// Format of `Address` is `ip:port`, or `ip:` to get random port.
type Server struct {
	Directory string
	Address   string
}

// ServeStatic starts a static server serving the contents of the `directory`
// over `address`. It returns the http.Server and random port number with error.
func (s *Server) ServeStatic(verbose bool) (*http.Server, int, error) {
	fs := http.FileServer(http.Dir(s.Directory))
	mux := http.NewServeMux()
	mux.Handle("/", fs)

	srv := &http.Server{Handler: mux}

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
		host := s.Address
		if host == "" {
			host = "127.0.0.1"
		} else {
			host = strings.Split(host, ":")[0]
		}
		log.Printf("Serving at http://%s:%d\n", host, p)
	}

	go func() {
		if err := srv.Serve(listener); err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()

	return srv, p, nil
}
