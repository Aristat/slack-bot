package commands

import (
	"fmt"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/db"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// LeavePollBlock command
type LeavePollBlock struct {
	Logger              *logrus.Logger
	ResponseWriter      http.ResponseWriter
	InteractionCallback *slack.InteractionCallback
	Client              *slack.Client
}

// Execute method
func (l *LeavePollBlock) Execute() {
	oldMessage := db.Messages[l.InteractionCallback.Channel.ID]
	if oldMessage == nil {
		return
	}

	var textBlock string
	index := findUser(oldMessage.Users, l.InteractionCallback.User.ID)
	if index == -1 {
		return
	}

	copy(oldMessage.Users[index:], oldMessage.Users[index+1:])
	oldMessage.Users = oldMessage.Users[:len(oldMessage.Users)-1]
	if len(oldMessage.Users) != 0 {
		for _, user := range oldMessage.Users {
			textBlock += fmt.Sprintf("<@%s>\n", user.ID)
		}
	} else {
		textBlock = "Poll is empty"
	}

	blocks, err := generatePollMessageBlocks(textBlock)
	if err != nil {
		l.Logger.Error(err.Error())
	}

	msgOptionBlock := slack.MsgOptionBlocks(blocks...)
	_, _, err = l.Client.PostMessage(
		l.InteractionCallback.Channel.ID,
		slack.MsgOptionReplaceOriginal(l.InteractionCallback.ResponseURL),
		msgOptionBlock,
	)
	if err != nil {
		l.Logger.Error(err.Error())
	}

	l.Logger.WithFields(logrus.Fields{
		"oldMessage": oldMessage,
	}).Info("leave-poll-block")
}
