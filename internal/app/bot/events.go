package bot

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/commands"
	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

// EventsHandler generate slack hanlder, /event-subscriptions
func (s *Slack) EventsHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sv, err := slack.NewSecretsVerifier(r.Header, s.SigningSecret)
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

	switch eventsAPIEvent.Type {
	case slackevents.URLVerification:
		var challengeR *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &challengeR)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text")
		w.Write([]byte(challengeR.Challenge))
	case slackevents.CallbackEvent:
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slackevents.MessageEvent:
			s.logger.WithFields(logrus.Fields{
				"inner_event": ev,
			}).Info()

			if ev.BotID == "" && ev.User == "" {
				break
			}

			if ev.BotID != "" {
				break
			}

			switch ev.Text {
			case utils.BeerMessageEvent:
				cmd := commands.BeerEvent{SlackEvent: &eventsAPIEvent, MessageEvent: ev, ResponseWriter: w, Client: s.Client, Logger: s.logger}
				cmd.Execute()
			}
		}
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
