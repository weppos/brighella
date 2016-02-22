package main

import (
	"fmt"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/miekg/dns"
)

const what = "brighella"
const dnsPrefix = "_frame"
const resolverAddress = "8.8.8.8"
const resolverPort = "53"

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
// Any request that is not for the root path / is automatically redirected
// to the root with a 302 status code. Only a request to / will enable the iframe.
func (s *Server) Root(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		s.TemporaryRedirect(w, r, "/")
	} else {
		targetURL, err := queryRedirectTarget(r.Host)
		// An error happened. For now, do not display the full error.
		if err != nil {
			http.Error(w, "Unable to find redirect target", http.StatusBadRequest)
			return
		}
		s.MaskedRedirect(w, r, targetURL)
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

func queryRedirectTarget(host string) (string, error) {
	targetFqdn := fmt.Sprintf("%s.%s", dnsPrefix, host)

	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(targetFqdn), dns.TypeTXT)
	r, _, err := c.Exchange(m, net.JoinHostPort(resolverAddress, resolverPort))

	if err != nil {
		log.Printf("[%s] Error querying %s: %v", host, targetFqdn, err)
		return "", err
	}

	if r.Rcode != dns.RcodeSuccess {
		err = fmt.Errorf("answer from %s not successful: %v", targetFqdn, dns.RcodeToString[r.Rcode])
		log.Printf("[%s] Error %s", host, err)
		return "", err
	}

	for _, a := range r.Answer {
		switch rr := a.(type) {
		case *dns.TXT:
			log.Printf("[%s] Found redirect target at %s: %v", host, targetFqdn, rr.Txt[0])
			return rr.Txt[0], nil
		}
	}

	err = fmt.Errorf("redirect target not found at %s", targetFqdn)
	log.Printf("[%s] Error %s", host, err)

	return "", err
}

type frame struct {
	Src string
}
