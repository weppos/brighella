package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

const what = "brighella"

var (
	httpPort string
)

func init() {
	httpPort = os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "5000"
	}
}

func main() {
	log.Printf("Starting %s...\n", what)

	server := NewServer()

	log.Printf("%s listening on %s...\n", what, httpPort)
	if err := http.ListenAndServe(":"+httpPort, server); err != nil {
		log.Panic(err)
	}
}

// Server represents a front-end web server.
type Server struct {
	// Router which handles incoming requests
	mux *http.ServeMux
}

// NewServer returns a new front-end web server that handles HTTP requests for the app.
func NewServer() *Server {
	router := http.NewServeMux()
	server := &Server{mux: router}
	router.HandleFunc("/", server.Root)
	return server
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// Root is the handler for the HTTP requests to /.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")
	fmt.Fprintln(w, `
<frameset rows="100%">
  <frameset cols="100%">
    <frame src="http://example.com/page.html" frameborder="0" scrolling="no">
  </frameset>
</frameset>
	`)
}
