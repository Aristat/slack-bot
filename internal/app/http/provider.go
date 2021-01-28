package http

import (
	"github.com/aristat/slack-bot/internal/app/router"
	"github.com/go-chi/chi"
	"github.com/google/wire"
	"github.com/sirupsen/logrus"
)

// Provider http server
func Provider(mux *chi.Mux, router *router.Router, logger *logrus.Logger) (*HTTP, func(), error) {
	g := New(mux, logger)
	router.Run(mux)
	return g, func() {}, nil
}

var (
	// ProviderProductionSet http server set
	ProviderProductionSet = wire.NewSet(Provider, Mux)
)
