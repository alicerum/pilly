package alarm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alicerum/pilly/pkg/db"
	alarmsModel "github.com/alicerum/pilly/pkg/db/daily/alarms"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	response = "Alarm successfully created"

	errorPersist = "Could not save alarm. Please try again."

	errorBadTime = "Could not parse time.\n" +
		"Make sure time is correct and is set up in a `hh:mm` format,\n" +
		"with hours being in a 24 hours format.\n" +
		"i.e.: `13:46`"
)

type Processor struct {
	db *db.Svc
}

func NewProcessor(db *db.Svc) *Processor {
	return &Processor{
		db: db,
	}
}

func createResult(msg *tgbotapi.Message, text string) *tgbotapi.MessageConfig {
	res := tgbotapi.NewMessage(msg.From.ID, text)
	return &res
}

func badTimeResult(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	return createResult(msg, errorBadTime)
}

func minutesToTime(minutes int) string {
	return fmt.Sprintf("%d:%02d", minutes/60, minutes%60)
}

func (p *Processor) Process(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	time := strings.Trim(msg.CommandArguments(), " ")
	parts := strings.Split(time, ":")
	if len(parts) != 2 {
		return badTimeResult(msg)
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Error().Err(err).Msg("could not parse hours in alarm command")
		return badTimeResult(msg)
	}
	if hours >= 24 {
		return badTimeResult(msg)
	}

	minutes, err := strconv.Atoi(parts[1])
	if err != nil {
		log.Error().Err(err).Msg("could not parse minutes in alarm command")
		return badTimeResult(msg)
	}
	if minutes >= 60 {
		return badTimeResult(msg)
	}

	alarm := alarmsModel.Alarm{
		UserID:  msg.From.ID,
		Minutes: hours*60 + minutes,
	}

	err = p.db.Alarms().Persist(&alarm)
	if err != nil {
		log.Error().Err(err).
			Int64("userId", msg.From.ID).
			Str("args", msg.CommandArguments()).
			Msg("could not persist alarm")
		return createResult(msg, errorPersist)
	}

	return createResult(msg, response)
}
