package main

import (
	"html/template"
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
	if r.URL.Path != "/" {
		s.TemporaryRedirect(w, r, "/")
	} else {
		s.MaskedRedirect(w, r, "http://example.com/page.html")
	}
}

func (s *Server) TemporaryRedirect(w http.ResponseWriter, r *http.Request, strURL string) {
	http.Redirect(w, r, strURL, http.StatusTemporaryRedirect)
}

func (s *Server) MaskedRedirect(w http.ResponseWriter, r *http.Request, strURL string) {
	w.Header().Set("Content-type", "text/html")

	t, _ := template.ParseFiles("redirect.tmpl")
	t.Execute(w, &frame{Src: strURL})
}

type frame struct {
	Src string
}
