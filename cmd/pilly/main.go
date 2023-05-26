package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alicerum/pilly/pkg/config"
	"github.com/alicerum/pilly/pkg/db/users"
	"github.com/rs/zerolog/log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func printErrAndExit(errMsg string, err error) {
	fmt.Fprintf(os.Stderr, "%v: %v\n", errMsg, err)
	os.Exit(1)
}

func main() {
	configPath := flag.String("config", "", "config path")
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		printErrAndExit("error while reading config file", err)
	}

	if err := cfg.Log.SetLogger(); err != nil {
		printErrAndExit("could not set logging level", err)
	}

	db, err := cfg.DB.CreateDB(ctx)
	if err != nil {
		printErrAndExit("could not connect to the db", err)
	}

	usersSvc := users.NewSvc(db)

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		printErrAndExit("could not create bot api", err)
	}

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	updates := bot.GetUpdatesChan(upd)
	for u := range updates {
		if u.Message != nil {
			log.Info().
				Str("from", u.Message.From.UserName).
				Str("text", u.Message.Text).
				Msg("got message")
		}
		u := users.User{
			ID:        u.Message.From.ID,
			Username:  u.Message.From.UserName,
			FirstName: u.Message.From.FirstName,
			LastName:  u.Message.From.LastName,
			Created:   time.Now().UTC(),
		}
		usersSvc.Persist(&u)
	}
}
