package router

import (
	"net/http"

	"github.com/aristat/slack-bot/internal/app/bot"
	"github.com/go-chi/chi"
	"github.com/slack-go/slack"
)

// Router struct
type Router struct {
	bot           *bot.Slack
	client        *slack.Client
	signingSecret string
}

// New router
func New(bot *bot.Slack) *Router {
	router := &Router{bot: bot, client: bot.Client, signingSecret: bot.SigningSecret}
	return router
}

// Run attach
func (router *Router) Run(chiRouter chi.Router) {
	chiRouter.HandleFunc("/slack/events", router.eventsHandler)
	chiRouter.HandleFunc("/slack/redirect", router.redirectHandler)
	chiRouter.HandleFunc("/slack/interative-endpoint", router.interativeHandler)
	chiRouter.HandleFunc("/slack/command", router.commandHandler)
}

func (router *Router) eventsHandler(w http.ResponseWriter, r *http.Request) {
	router.bot.EventsHandler(w, r)
}

func (router *Router) interativeHandler(w http.ResponseWriter, r *http.Request) {
	router.bot.InterativeHandler(w, r)
}

func (router *Router) commandHandler(w http.ResponseWriter, r *http.Request) {
	router.bot.CommandHandler(w, r)
}

func (router *Router) redirectHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "text")
	w.Write([]byte("ok"))
}
