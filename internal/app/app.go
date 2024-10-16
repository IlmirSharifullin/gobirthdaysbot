package app

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/handlers"
	"telegram-bot/internal/app/states"
)

func initStates() {
	states.AnyState.Handler = handlers.StartCommand
	states.BirthdayNameState.Handler = handlers.BirthdayName
	states.BirthdayDateState.Handler = handlers.BirthdayDate
	states.BirthdayAdditionalState.Handler = handlers.BirthdayAdditional
}

func GetUpdates(updates tg.UpdatesChannel, ctx common.Context) {
	initStates()
	for update := range updates {
		var err error

		if !common.MsgNotNil(update) {
			continue
		}
		if update.Message.Command() == "start" {
			err = handlers.StartCommand(ctx, update)
		} else if update.Message.Command() == "add" {
			err = handlers.BirthdayAddCommand(ctx, update)
		} else {
			err = ctx.Serve(update)
		}

		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("%v", err))
		}
	}
}
