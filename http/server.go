// Package http provides an HTTP server with a configurable server and router.
package http

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

// Server carries an HTTP server and a router that can be configured through Option when instantiating it.
type Server struct {
	server *http.Server
	router *mux.Router
}

// NewServer creates a new HTTP Server with the given listening address and configured with the provided Option.
func NewServer(addr string, opts ...Option) *Server {
	router := mux.NewRouter()
	s := &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           router,
			ReadHeaderTimeout: 5 * time.Second,
		},
		router: router,
	}

	WithOptions(opts...)(s)

	return s
}

// Option is a function to configure a Server.
type Option func(*Server)

// WithOptions construct an Option that applies multiple Option to a Server.
func WithOptions(opts ...Option) Option {
	return func(s *Server) {
		for _, opt := range opts {
			opt(s)
		}
	}
}

// WithRouterOption construct an Option that configure the Server's router.
func WithRouterOption(opt func(*mux.Router)) Option {
	return func(s *Server) {
		opt(s.router)
	}
}

// WithServerOption construct an Option that configure the Server's http.Server.
func WithServerOption(opt func(*http.Server)) Option {
	return func(s *Server) {
		opt(s.server)
	}
}

// WithRoute construct an Option that adds a route to the Server's router.
func WithRoute(method, path string, handler http.Handler) Option {
	return WithRouterOption(func(r *mux.Router) {
		r.Methods(method).Path(path).Handler(handler)
	})
}

// Listen runs the server in a blocking way. In returns either if an error occur which in that case returns the error,
// or if the server is stopped by a signal (i.e. `SIGINT` or `SIGTERM`).
func (s *Server) Listen() error {
	listenErr := make(chan error, 1)
	go func() {
		listenErr <- s.server.ListenAndServe()
	}()

	kill := make(chan os.Signal, 1)
	signal.Notify(kill, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-listenErr:
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		case <-kill:
			if err := s.server.Shutdown(context.Background()); err != nil {
				return err
			}
		}
	}
}
