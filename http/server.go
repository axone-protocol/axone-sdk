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

type Server struct {
	server *http.Server
	router *mux.Router
}

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
