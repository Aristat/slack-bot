// +build wireinject

package http

import (
	"github.com/aristat/slack-bot/internal/app/logger"
	"github.com/aristat/slack-bot/internal/app/router"
	"github.com/google/wire"
)

// Build
func Build() (*HTTP, func(), error) {
	panic(wire.Build(
		ProviderProductionSet,
		router.ProviderProductionSet,
		logger.ProviderProductionSet,
	))
}
