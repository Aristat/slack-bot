package commands

import (
	"github.com/aristat/slack-bot/internal/app/db"
	"github.com/aristat/slack-bot/internal/app/utils"
	"github.com/slack-go/slack"
)

func generatePollMessageBlocks(textBlock string) ([]slack.Block, error) {
	headerBlock := slack.HeaderBlock{
		Type: "header",
		Text: &slack.TextBlockObject{Type: "plain_text", Text: "Poll"},
	}

	contextBlock := slack.NewContextBlock("",
		&slack.TextBlockObject{
			Type: "mrkdwn",
			Text: textBlock,
		},
	)

	joinElement := slack.ButtonBlockElement{
		ActionID: utils.ActionPollJoin,
		Type:     "button",
		Style:    "primary",
		Text:     &slack.TextBlockObject{Type: "plain_text", Text: "Join"},
		Value:    utils.ActionPollJoin,
	}

	cancelElement := slack.ButtonBlockElement{
		ActionID: utils.ActionPollLeave,
		Type:     "button",
		Style:    "danger",
		Text:     &slack.TextBlockObject{Type: "plain_text", Text: "Leave"},
		Value:    utils.ActionPollLeave,
	}

	actionBlock := slack.ActionBlock{
		Type:     "actions",
		BlockID:  utils.PollBlockID,
		Elements: &slack.BlockElements{ElementSet: []slack.BlockElement{&joinElement, &cancelElement}},
	}

	return []slack.Block{headerBlock, contextBlock, actionBlock}, nil
}

func findUser(users []db.User, userID string) (index int) {
	for i, user := range users {
		if user.ID == userID {
			return i
		}
	}

	return -1
}
