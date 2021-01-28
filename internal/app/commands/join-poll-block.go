package commands

import (
	"fmt"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/db"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// JoinPollBlock command
type JoinPollBlock struct {
	Logger              *logrus.Logger
	ResponseWriter      http.ResponseWriter
	InteractionCallback *slack.InteractionCallback
	Client              *slack.Client
}

// Execute method
func (j *JoinPollBlock) Execute() {
	var textBlock string
	var users []db.User

	newUser := db.User{ID: j.InteractionCallback.User.ID}
	oldMessage := db.Messages[j.InteractionCallback.Channel.ID]
	if oldMessage != nil {
		for _, user := range oldMessage.Users {
			textBlock += fmt.Sprintf("<@%s>\n", user.ID)
		}

		// skip check that User already exist(for testing)
		users = append(oldMessage.Users, newUser)
		textBlock += fmt.Sprintf("<@%s>\n", newUser.ID)
	} else {
		return
	}

	blocks, err := generatePollMessageBlocks(textBlock)
	if err != nil {
		j.Logger.Error(err.Error())
	}

	msgOptionBlock := slack.MsgOptionBlocks(blocks...)
	respChannel, respTimestamp, err := j.Client.PostMessage(
		j.InteractionCallback.Channel.ID,
		slack.MsgOptionReplaceOriginal(j.InteractionCallback.ResponseURL),
		msgOptionBlock,
	)
	if err != nil {
		j.Logger.Error(err.Error())
	}

	j.Logger.WithFields(logrus.Fields{
		"respChannel":   respChannel,
		"respTimestamp": respTimestamp,
		"oldMessage":    oldMessage,
	}).Info("join-poll-block")

	db.Messages[j.InteractionCallback.Channel.ID] = &db.Message{
		Channel:   oldMessage.Channel,
		Timestamp: oldMessage.Timestamp,
		Users:     users,
	}
}
