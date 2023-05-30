package alarms

import (
	"fmt"

	"github.com/alicerum/pilly/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rs/zerolog/log"
)

const (
	errSelectAlarms = "Error while looking for alarms"
	errNoAlarms     = "No alarms exist for this user"

	response = "These alarms are set up for the current user.\n" +
		"Tap on the alarm to delete it."

	buttonDataFormat = "alarmdel %d"
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

func minutesToTime(minutes int) string {
	return fmt.Sprintf("%d:%02d", minutes/60, minutes%60)
}

func (p *Processor) Process(msg *tgbotapi.Message) *tgbotapi.MessageConfig {
	alarms, err := p.db.Alarms().GetByUser(msg.From.ID)
	if err != nil {
		log.Error().Err(err).Msg("error while selecting alarms")
		return createResult(msg, errSelectAlarms)
	}

	if len(alarms) == 0 {
		return createResult(msg, errNoAlarms)
	}

	responseMsg := createResult(msg, response)
	var keyboard [][]tgbotapi.InlineKeyboardButton
	var currentRow []tgbotapi.InlineKeyboardButton
	for i, alarm := range alarms {
		currentRow = append(
			currentRow,
			tgbotapi.NewInlineKeyboardButtonData(
				minutesToTime(alarm.Minutes),
				fmt.Sprintf(buttonDataFormat, alarm.ID),
			),
		)
		if len(currentRow) == 3 || i == len(alarms)-1 {
			keyboard = append(keyboard, currentRow)
			currentRow = make([]tgbotapi.InlineKeyboardButton, 0, 3)
		}
	}

	markup := tgbotapi.NewInlineKeyboardMarkup(keyboard...)
	responseMsg.ReplyMarkup = markup

	return responseMsg
}
