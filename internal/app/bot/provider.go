package bot

import (
	"log"
	"os"

	"github.com/google/wire"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// Slack struct
type Slack struct {
	Client        *slack.Client
	logger        *logrus.Logger
	SigningSecret string
}

// Provider initializer
func Provider(logger *logrus.Logger) (*Slack, func(), error) {
	client := slack.New(
		os.Getenv("SLACK_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	signingSecret := os.Getenv("SLACK_SIGNING_SECRET")

	config := &Slack{Client: client, SigningSecret: signingSecret, logger: logger}
	return config, func() {}, nil
}

var (
	// ProviderProductionSet wire set
	ProviderProductionSet = wire.NewSet(Provider)
)
