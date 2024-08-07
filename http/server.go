package http

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	server *http.Server
	router *mux.Router
}

func NewServer(addr string, opts ...Option) *Server {
	router := mux.NewRouter()
	s := &Server{
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
		router: router,
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Option func(*Server)

func WithOptions(opts ...Option) Option {
	return func(s *Server) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

func WithRoute(method, path string, handler http.Handler) Option {
	return func(s *Server) {
		s.router.Methods(method).Path(path).Handler(handler)
	}
}
