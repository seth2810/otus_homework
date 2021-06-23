package internalhttp

import (
	"context"
	"io"
	"net/http"
)

type Server struct {
	logger Logger
	server *http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

type HelloHandler struct{}

func NewServer(address string, logger Logger, app Application) *Server {
	mux := http.NewServeMux()
	mux.Handle("/hello", loggingMiddleware(&HelloHandler{}, logger))

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return &Server{logger, server}
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}

	<-ctx.Done()

	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, "ok")
}
