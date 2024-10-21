package app

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/robfig/cron/v3"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/handlers"
	"telegram-bot/internal/app/states"
	mycron "telegram-bot/internal/cron"
	"telegram-bot/internal/storage"
)

func initStates() {
	states.AnyState.Handler = handlers.StartCommand
	states.BirthdayNameState.Handler = handlers.BirthdayName
	states.BirthdayDateState.Handler = handlers.BirthdayDate
	states.BirthdayAdditionalState.Handler = handlers.BirthdayAdditional
}

var notificationDays = storage.NotificationDays{
	WeekBefore:      true,
	ThreeDaysBefore: false,
	DayBefore:       true,
	AtBirthDay:      true,
}

func GetUpdates(updates tg.UpdatesChannel, ctx common.Context) {
	initStates()

	c := cron.New()
	_, err := c.AddFunc("0 10 * * *", func() { mycron.ServeBirthdaysNotifications(ctx, notificationDays) })

	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("%s", err))
	}
	c.Start()

	for update := range updates {
		var err error

		if !common.MsgNotNil(update) {
			continue
		}
		if update.Message.Command() == "start" {
			err = handlers.StartCommand(ctx, update)
		} else if update.Message.Command() == "add" || update.Message.Text == "Add a birthday" {
			err = handlers.BirthdayAddCommand(ctx, update)
		} else if update.Message.Command() == "get" || update.Message.Text == "Get all birthdays" {
			err = handlers.GetAllBirthdays(ctx, update)
		} else {
			err = ctx.Serve(update)
		}

		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("%v", err))
		}
	}
}
