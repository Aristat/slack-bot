package main

import (
	"errors"
	"log"
	"os"

	srv "github.com/aristat/slack-bot/internal/app/http"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Commands: []*cli.Command{
			{
				Name:        "start",
				Aliases:     []string{"s"},
				Usage:       "run slack bot server",
				Description: "slack bot",
				Flags: []cli.Flag{
					&cli.StringFlag{Name: "envFile", Value: ".env.development"},
				},
				Action: func(c *cli.Context) error {
					fileName := c.String("envFile")

					err := godotenv.Load(fileName)
					if err != nil {
						return errors.New("Error loading .env file")
					}

					srv, cleanup, err := srv.Build()
					defer cleanup()
					if err != nil {
						return err
					}

					err = srv.Serve().ListenAndServe()
					return err
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal("cli error ", err)
		os.Exit(1)
	}
}
