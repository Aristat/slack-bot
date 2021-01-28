package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

// HTTP struct
type HTTP struct {
	mux    *chi.Mux
	logger *logrus.Logger
}

// Serve create
func (s *HTTP) Serve() (server *http.Server) {
	s.logger.Info("Initialize http server")

	server = &http.Server{
		Addr:         "127.0.0.1:3000",
		Handler:      s.mux,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server
}

// New http struct
func New(mux *chi.Mux, logger *logrus.Logger) *HTTP {
	return &HTTP{
		mux:    mux,
		logger: logger,
	}
}
