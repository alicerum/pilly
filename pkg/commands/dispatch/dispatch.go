package dispatch

import (
	"errors"
	"strings"

	"github.com/alicerum/pilly/pkg/commands/alarms/alarm"
	"github.com/alicerum/pilly/pkg/commands/alarms/alarmdel"
	"github.com/alicerum/pilly/pkg/commands/alarms/alarms"
	"github.com/alicerum/pilly/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Cmd string

const (
	Alarms Cmd = "alarms"
	Alarm  Cmd = "alarm"

	Alarmdel Cmd = "alarmdel"
)

var (
	ErrNotCmd               = errors.New("not a command")
	ErrCmdNotSupported      = errors.New("command not supported")
	ErrCallbackNotSupported = errors.New("callback not supported")
)

type CmdProcessor interface {
	Process(msg *tgbotapi.Message) *tgbotapi.MessageConfig
}

type CallbackProcessor interface {
	Process(
		query *tgbotapi.CallbackQuery,
		args string,
	) (*tgbotapi.CallbackConfig, *tgbotapi.MessageConfig)
}

type Dispatcher struct {
	cmdProcessors      map[Cmd]CmdProcessor
	callbackProcessors map[Cmd]CallbackProcessor
}

func fillCmdProcessors(db *db.Svc) map[Cmd]CmdProcessor {
	return map[Cmd]CmdProcessor{
		Alarms: alarms.NewProcessor(db),
		Alarm:  alarm.NewProcessor(db),
	}
}

func fillCallbackProcessors(db *db.Svc) map[Cmd]CallbackProcessor {
	return map[Cmd]CallbackProcessor{
		Alarmdel: alarmdel.NewProcessor(db),
	}
}

func NewDispatcher(db *db.Svc) *Dispatcher {
	return &Dispatcher{
		cmdProcessors:      fillCmdProcessors(db),
		callbackProcessors: fillCallbackProcessors(db),
	}
}

func (d *Dispatcher) ProcessCmd(
	msg *tgbotapi.Message,
) (*tgbotapi.MessageConfig, error) {
	if !msg.IsCommand() {
		return nil, ErrNotCmd
	}

	proc, ok := d.cmdProcessors[Cmd(msg.Command())]
	if !ok {
		return nil, ErrCmdNotSupported
	}

	return proc.Process(msg), nil
}

func (d *Dispatcher) ProcessCallback(
	query *tgbotapi.CallbackQuery,
) (*tgbotapi.CallbackConfig, *tgbotapi.MessageConfig, error) {

	parts := strings.Split(query.Data, " ")
	proc, ok := d.callbackProcessors[Cmd(parts[0])]
	if !ok {
		return nil, nil, ErrCmdNotSupported
	}

	cbk, msg := proc.Process(query, strings.Join(parts[1:], " "))
	return cbk, msg, nil
}
