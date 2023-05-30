package alarmdel

import (
	"strconv"

	"github.com/alicerum/pilly/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	errCouldNotDelete = "Could not delete alarm"
	success           = "Alarm has been successfully deleted"
)

type Processor struct {
	db *db.Svc
}

func NewProcessor(db *db.Svc) *Processor {
	return &Processor{
		db: db,
	}
}

func creatCbk(query *tgbotapi.CallbackQuery, text string) *tgbotapi.CallbackConfig {
	callback := tgbotapi.NewCallback(query.ID, text)
	return &callback
}

func createMsg(query *tgbotapi.CallbackQuery, text string) *tgbotapi.MessageConfig {
	res := tgbotapi.NewMessage(query.Message.Chat.ID, text)
	return &res
}

func createCbkAndMsg(
	query *tgbotapi.CallbackQuery,
	text string,
) (*tgbotapi.CallbackConfig, *tgbotapi.MessageConfig) {
	callback := creatCbk(query, text)
	msg := createMsg(query, text)
	return callback, msg
}

func (p *Processor) Process(
	query *tgbotapi.CallbackQuery,
	args string,
) (*tgbotapi.CallbackConfig, *tgbotapi.MessageConfig) {
	alarmID, err := strconv.Atoi(args)
	if err != nil {
		log.Error().Err(err).Str("args", args).Msg("could not parse alarmdel arg")
		return createCbkAndMsg(query, errCouldNotDelete)
	}

	err = p.db.Alarms().DeleteByID(alarmID)
	if err != nil {
		log.Error().Err(err).Int("alarmID", alarmID).Msg("could not delete alarm from db")
		return createCbkAndMsg(query, errCouldNotDelete)
	}

	return createCbkAndMsg(query, success)
}
