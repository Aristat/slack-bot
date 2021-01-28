package commands

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// BeerCallback command
type BeerCallback struct {
	Logger              *logrus.Logger
	ResponseWriter      http.ResponseWriter
	InteractionCallback *slack.InteractionCallback
	Client              *slack.Client
}

// Execute method
func (b *BeerCallback) Execute() {
	action := b.InteractionCallback.ActionCallback.AttachmentActions[0]

	switch action.Name {
	case utils.ActionBeerSelect:
		value := action.SelectedOptions[0].Value

		// Overwrite original drop down message.
		originalMessage := b.InteractionCallback.OriginalMessage
		originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
		originalMessage.ReplaceOriginal = true
		originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
			{
				Name:  utils.ActionBeerStart,
				Text:  "Yes",
				Type:  "button",
				Value: utils.ActionBeerStart,
				Style: "primary",
			},
			{
				Name:  utils.ActionBeerCancel,
				Text:  "No",
				Type:  "button",
				Value: utils.ActionBeerCancel,
				Style: "danger",
			},
		}

		b.ResponseWriter.Header().Add("Content-type", "application/json")
		b.ResponseWriter.WriteHeader(http.StatusOK)
		json.NewEncoder(b.ResponseWriter).Encode(&originalMessage)
		return
	case utils.ActionBeerStart:
		title := ":ok: your order was submitted! yay!"

		msgRef := slack.NewRefToMessage(b.InteractionCallback.Channel.ID, b.InteractionCallback.OriginalMessage.Timestamp)
		listPins, _, _ := b.Client.ListPins(b.InteractionCallback.Channel.ID)

		b.Logger.WithFields(logrus.Fields{
			"list_pins": listPins,
		}).Info()
		if err := b.Client.RemovePin(b.InteractionCallback.Channel.ID, msgRef); err != nil {
			fmt.Printf("Error adding pin: %s\n", err)
		}

		beerResponseMessage(b.ResponseWriter, b.InteractionCallback.OriginalMessage, title, "")
		return
	case utils.ActionBeerCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", b.InteractionCallback.User.Name)
		beerResponseMessage(b.ResponseWriter, b.InteractionCallback.OriginalMessage, title, "")
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		b.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func beerResponseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.ReplaceOriginal = true
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
