package commands

import (
	"fmt"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/db"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// PollCommand command
type PollCommand struct {
	Logger            *logrus.Logger
	ResponseWriter    http.ResponseWriter
	SlackSlashCommand *slack.SlashCommand
	Client            *slack.Client
}

// Execute method
func (pCommand *PollCommand) Execute() {
	var textBlock string
	var users []db.User

	oldMessage := db.Messages[pCommand.SlackSlashCommand.ChannelID]
	if oldMessage != nil {
		_, _, err := pCommand.Client.DeleteMessage(oldMessage.Channel, oldMessage.Timestamp)
		if err != nil {
			pCommand.Logger.Error(err.Error())
		}

		for _, user := range oldMessage.Users {
			textBlock += fmt.Sprintf("<@%s>\n", user.ID)
		}

		users = oldMessage.Users
	} else {
		textBlock = "Poll is empty"
		users = []db.User{}
	}

	blocks, err := generatePollMessageBlocks(textBlock)
	if err != nil {
		pCommand.Logger.Error(err.Error())
	}

	msgOptionBlock := slack.MsgOptionBlocks(blocks...)
	respChannel, respTimestamp, err := pCommand.Client.PostMessage(
		pCommand.SlackSlashCommand.ChannelID,
		msgOptionBlock,
	)

	db.Messages[respChannel] = &db.Message{Channel: respChannel, Timestamp: respTimestamp, Users: users}
}
