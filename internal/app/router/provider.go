package router

import (
	"github.com/aristat/slack-bot/internal/app/bot"
	"github.com/google/wire"
)

// Provider slack config
func Provider(bot *bot.Slack) (*Router, func(), error) {
	g := New(bot)
	return g, func() {}, nil
}

var (
	// ProviderProductionSet slack config set
	ProviderProductionSet = wire.NewSet(Provider, bot.ProviderProductionSet)
)
