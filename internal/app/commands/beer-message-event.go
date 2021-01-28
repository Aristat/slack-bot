package commands

import (
	"net/http"

	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// BeerEvent command
type BeerEvent struct {
	Logger         *logrus.Logger
	ResponseWriter http.ResponseWriter
	SlackEvent     *slackevents.EventsAPIEvent
	Client         *slack.Client
	MessageEvent   *slackevents.MessageEvent
}

// Execute method
func (b *BeerEvent) Execute() {
	attachment := slack.Attachment{
		Text:       "Which beer do you want? :beer:",
		Color:      "#f9a41b",
		CallbackID: utils.BeerCallbackID,
		Actions: []slack.AttachmentAction{
			{
				Name: utils.ActionBeerSelect,
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "Asahi Super Dry",
						Value: "Asahi Super Dry",
					},
					{
						Text:  "Kirin Lager Beer",
						Value: "Kirin Lager Beer",
					},
					{
						Text:  "Sapporo Black Label",
						Value: "Sapporo Black Label",
					},
					{
						Text:  "Suntory Malts",
						Value: "Suntory Malts",
					},
					{
						Text:  "Yona Yona Ale",
						Value: "Yona Yona Ale",
					},
				},
				Value: utils.ActionBeerSelect,
			},

			{
				Name:  utils.ActionBeerCancel,
				Text:  "Cancel",
				Type:  "button",
				Style: "danger",
				Value: utils.ActionBeerCancel,
			},
		},
	}
	channelID, timestamp, err := b.Client.PostMessage(b.MessageEvent.Channel, slack.MsgOptionAttachments(attachment))
	if err != nil {
		b.Logger.Error(err.Error())
	}

	msgRef := slack.NewRefToMessage(channelID, timestamp)
	if err = b.Client.AddPin(channelID, msgRef); err != nil {
		b.Logger.Error(err.Error())
	}
}
