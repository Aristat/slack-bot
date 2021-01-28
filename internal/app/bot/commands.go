package bot

import (
	"io"
	"io/ioutil"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/commands"
	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// CommandHandler generate slack hanlder, /slash-commands
func (s *Slack) CommandHandler(w http.ResponseWriter, r *http.Request) {
	verifier, err := slack.NewSecretsVerifier(r.Header, s.SigningSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = ioutil.NopCloser(io.TeeReader(r.Body, &verifier))

	slackSlashCommand, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"slash_command": slackSlashCommand,
	}).Info()

	switch slackSlashCommand.Command {
	case utils.PollCommand:
		cmd := commands.PollCommand{SlackSlashCommand: &slackSlashCommand, ResponseWriter: w, Client: s.Client, Logger: s.logger}
		cmd.Execute()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
