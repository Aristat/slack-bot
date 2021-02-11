# Slack Bot

# Table of Contents

- [Overview](#overview)
- [Package list](#package-list)
- [Run bot](#run-bot)

# Overview

This is an example slack bot

Based on https://github.com/tcnksm/go-slack-interactive ("beer" logic) and simple poll

# Package list

Packages which use in this example project

1. [slack](https://github.com/slack-go/slack) - slack library
2. [wire](https://github.com/google/wire) - dependency Injection
3. [chi](https://github.com/go-chi/chi) - HTTP router
4. [godotenv](https://github.com/joho/godotenv) - load env variables from files
5. [logrus](https://github.com/sirupsen/logrus) - logger
6. [cli](https://github.com/urfave/cli) - simple and fast CLI

# Run bot

1. Install golang, packages
2. Add `.env.development` file with `SLACK_TOKEN` and `SLACK_SIGNING_SECRET` in repository
3. Run
    ```$xslt
    go run main.go start
    ```
4. Install ngrok
5. Setup in `api.slack.com/`:
    1. `Interactivity Request URL = https://#.ngrok.io/slack/interative-endpoint`
    2. `Slash Commands '/poll' = https://#.ngrok.io/slack/command`
    3. `OAuth & Permissions Redirect URLs = https://#.ngrok.io/slack/redirect` 
    4. `Scopes channels:history, chat:write, commands, pins:read, pins:write`
    5. `Event Subscriptions Request URL = https://#.ngrok.io/slack/events`
6. Add bot to app
7. Write `give beer` in a chat or run `/poll` command 
