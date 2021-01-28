package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const (
	// action is used for slack attament action.
	actionSelect = "select"
	actionStart  = "start"
	actionCancel = "cancel"
)

var (
	client        *slack.Client
	signingSecret string
)

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if _, err := sv.Write(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err := sv.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(r.Challenge))
	}
	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.AppMentionEvent:
			client.PostMessage(ev.Channel, slack.MsgOptionText("Yes, hello.", false))
		case *slackevents.MessageEvent:
			if ev.BotID == "" && ev.User == "" {
				break
			}

			if ev.BotID != "" {
				break
			}
			fmt.Printf("slackevents %+v\n", ev)

			client.PostMessage(ev.Channel, slack.MsgOptionText("Hello!", false))

			if ev.Text == "test" {
				attachment := slack.Attachment{
					Text:       "Which beer do you want? :beer:",
					Color:      "#f9a41b",
					CallbackID: "beer",
					Actions: []slack.AttachmentAction{
						{
							Name: actionSelect,
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
							Value: actionSelect,
						},

						{
							Name:  actionCancel,
							Text:  "Cancel",
							Type:  "button",
							Style: "danger",
							Value: actionCancel,
						},
					},
				}
				client.PostMessage(ev.Channel, slack.MsgOptionAttachments(attachment))
			}
		}
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	w.Write([]byte("ok"))
}

func interativeEndpointHandler(w http.ResponseWriter, r *http.Request) {
	var ev slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &ev)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Println("OriginalMessage1", ev.OriginalMessage)
	fmt.Printf("AttachmentActions %+v\n", ev.ActionCallback.AttachmentActions[0])

	action := ev.ActionCallback.AttachmentActions[0]
	switch action.Name {
	case actionSelect:
		value := action.SelectedOptions[0].Value

		// Overwrite original drop down message.
		originalMessage := ev.OriginalMessage

		fmt.Println("OriginalMessage2", originalMessage)

		originalMessage.Attachments[0].Text = fmt.Sprintf("OK to order %s ?", strings.Title(value))
		originalMessage.ReplaceOriginal = true
		originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
			{
				Name:  actionStart,
				Text:  "Yes",
				Type:  "button",
				Value: actionStart,
				Style: "primary",
			},
			{
				Name:  actionCancel,
				Text:  "No",
				Type:  "button",
				Value: actionCancel,
				Style: "danger",
			},
		}

		// client.UpdateMessage(originalMessage.Channel, originalMessage.Timestamp)

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&originalMessage)
		return
	case actionStart:
		title := ":ok: your order was submitted! yay!"
		responseMessage(w, ev.OriginalMessage, title, "")
		return
	case actionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", ev.User.Name)
		responseMessage(w, ev.OriginalMessage, title, "")
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	fmt.Printf("original message %+v\n", original)
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

func main() {
	err := godotenv.Load(".env.development")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client = slack.New(
		os.Getenv("SLACK_TOKEN"),
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)
	signingSecret = os.Getenv("SLACK_SIGNING_SECRET")

	router := mux.NewRouter()
	router.HandleFunc("/slack/events", eventsHandler)
	router.HandleFunc("/slack/redirect", redirectHandler)
	router.HandleFunc("/slack/interative-endpoint", interativeEndpointHandler)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:3000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	fmt.Print("Start server\n")
	log.Fatal(srv.ListenAndServe())
}
