package webserver

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type WebServer struct {
	WebServerPort string
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	Middlewares   []Middleware
}

type Middleware func(h http.Handler) http.Handler

func NewWebServer(port string) *WebServer {
	return &WebServer{
		WebServerPort: port,
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		Middlewares:   make([]Middleware, 0),
	}
}

func (s *WebServer) AddHandler(pattern string, handler http.HandlerFunc) {
	s.Handlers[pattern] = handler
}

func (s *WebServer) AddMiddleware(middleware Middleware) {
	s.Middlewares = append(s.Middlewares, middleware)
}

func (s *WebServer) Start() {

	for _, mid := range s.Middlewares {
		s.Router.Use(mid)
	}

	for pattern, handler := range s.Handlers {
		s.Router.Handle(pattern, handler)
	}

	server := http.Server{
		Addr:    s.WebServerPort,
		Handler: s.Router,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("Could not listen on %s, %v\n", server.Addr, err)
	}

}
