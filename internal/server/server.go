package server

import (
	"context"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log/slog"
	"net/http"
)

type Server struct {
	srv       *http.Server
	lg        *slog.Logger
	storage   *gorm.DB
	validator *validator.Validate
}

func New(lg *slog.Logger, validator *validator.Validate, port string, storage *gorm.DB) *Server {

	s := &Server{
		lg:        lg,
		storage:   storage,
		validator: validator,
	}
	lg.Info("Initializing HTTP server", "port", port)

	r := http.NewServeMux()
	r.HandleFunc("/questions/", s.HandleQuestion)
	r.HandleFunc("/answers/", s.HandleAnswer)

	lg.Info("Routes registered", "routes", []string{"/questions/", "/answers/"})

	s.srv = &http.Server{
		Addr:    port,
		Handler: r,
	}
	lg.Info("Server instance created")
	return s
}
func (s *Server) Run() error {
	s.lg.Info("server running", "port", s.srv.Addr)
	return s.srv.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	s.lg.Info("server shutting down")
	err := s.srv.Shutdown(ctx)
	if err != nil {
		s.lg.Error("error for start shutdown")
	}
	return nil
}
