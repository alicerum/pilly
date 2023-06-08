package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/alicerum/pilly/pkg/commands/dispatch"
	"github.com/alicerum/pilly/pkg/config"
	pillyDB "github.com/alicerum/pilly/pkg/db"
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

	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		printErrAndExit("error while reading config file", err)
	}

	fmt.Printf("token is %v\n", cfg.Bot.Token)

	if err := cfg.Log.SetLogger(); err != nil {
		printErrAndExit("could not set logging level", err)
	}

	db, err := cfg.DB.CreateDB(ctx)
	if err != nil {
		printErrAndExit("could not connect to the db", err)
	}
	defer db.Close()

	dbSvc := pillyDB.New(db)
	dispatcher := dispatch.NewDispatcher(dbSvc)

	bot, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		printErrAndExit("could not create bot api", err)
	}

	upd := tgbotapi.NewUpdate(0)
	upd.Timeout = 60

	updates := bot.GetUpdatesChan(upd)
	for update := range updates {
		// ignore all public communications
		if !update.FromChat().IsPrivate() {
			continue
		}

		if update.Message != nil {
			user := users.User{
				ID:        update.Message.From.ID,
				Username:  update.Message.From.UserName,
				FirstName: update.Message.From.FirstName,
				LastName:  update.Message.From.LastName,
				Created:   time.Now().UTC(),
			}
			dbSvc.Users().Persist(&user)

			if update.Message.IsCommand() {
				response, err := dispatcher.ProcessCmd(update.Message)
				if err != nil {
					log.Error().Err(err).Msg("could not process tg message")
					continue
				}

				_, err = bot.Send(response)
				if err != nil {
					log.Error().Err(err).Msg("could not send response")
					continue
				}
			} else {
				msg, err := dispatcher.ProcessInput(update.Message)
				if err != nil {
					log.Error().Err(err).Msg("error while processing user input")
					continue
				}
				if msg == nil {
					log.Debug().
						Int64("userID", update.Message.From.ID).
						Msg("not expected input from user")
					continue
				}

				if _, err := bot.Send(msg); err != nil {
					log.Error().Err(err).Msg("could not send input response to user")
				}
			}
		}

		if update.CallbackQuery != nil {
			cbk, msg, err := dispatcher.ProcessCallback(update.CallbackQuery)
			if err != nil {
				log.Error().Err(err).Msg("could not process callback")
				continue
			}

			if _, err := bot.Request(cbk); err != nil {
				log.Error().Err(err).Msg("could not request callback")
				continue
			}

			if _, err := bot.Send(msg); err != nil {
				log.Error().Err(err).Msg("could not send callback response")
			}
		}
	}
}
