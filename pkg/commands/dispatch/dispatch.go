package dispatch

import (
	"errors"

	"github.com/alicerum/pilly/pkg/commands/alarms"
	"github.com/alicerum/pilly/pkg/db"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Cmd string

const (
	Alarms Cmd = "alarms"
)

var (
	ErrNotCmd          = errors.New("not a command")
	ErrCmdNotSupported = errors.New("command not supported")
)

type CmdProcessor interface {
	Process(msg *tgbotapi.Message) *tgbotapi.MessageConfig
}

type Dispatcher struct {
	processors map[Cmd]CmdProcessor
}

func fillCmdProcessors(db *db.Svc) map[Cmd]CmdProcessor {
	return map[Cmd]CmdProcessor{
		Alarms: alarms.NewProcessor(db),
	}
}

func NewDispatcher(db *db.Svc) *Dispatcher {
	return &Dispatcher{
		processors: fillCmdProcessors(db),
	}
}

func (d *Dispatcher) Process(
	msg *tgbotapi.Message,
) (*tgbotapi.MessageConfig, error) {
	if !msg.IsCommand() {
		return nil, ErrNotCmd
	}

	proc, ok := d.processors[Cmd(msg.Command())]
	if !ok {
		return nil, ErrCmdNotSupported
	}

	return proc.Process(msg), nil
}
