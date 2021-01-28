package bot

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aristat/slack-bot/internal/app/commands"
	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/sirupsen/logrus"
	"github.com/slack-go/slack"
)

// InterativeHandler generate slack hanlder, /interactive-messages
func (s *Slack) InterativeHandler(w http.ResponseWriter, r *http.Request) {
	var interactionCallback slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &interactionCallback)
	if err != nil {
		fmt.Printf(err.Error())
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	s.logger.WithFields(logrus.Fields{
		"callback": interactionCallback,
	}).Info()

	if len(interactionCallback.ActionCallback.AttachmentActions) != 0 {
		switch interactionCallback.CallbackID {
		case utils.BeerCallbackID:
			cmd := commands.BeerCallback{InteractionCallback: &interactionCallback, ResponseWriter: w, Client: s.Client, Logger: s.logger}
			cmd.Execute()
		}
	}

	if len(interactionCallback.ActionCallback.BlockActions) != 0 {
		action := interactionCallback.ActionCallback.BlockActions[0]

		switch action.ActionID {
		case utils.ActionPollJoin:
			cmd := commands.JoinPollBlock{InteractionCallback: &interactionCallback, ResponseWriter: w, Client: s.Client, Logger: s.logger}
			cmd.Execute()
		case utils.ActionPollLeave:
			cmd := commands.LeavePollBlock{InteractionCallback: &interactionCallback, ResponseWriter: w, Client: s.Client, Logger: s.logger}
			cmd.Execute()
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
