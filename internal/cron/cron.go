package mycron

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/storage"
)

func ServeBirthdaysNotifications(ctx common.Context, nd storage.NotificationDays) {
	ctx.Logger().Info("scheduler: start job")
	birthdays, err := ctx.Db().GetFilteredBirthdays(nd)

	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("scheduler: error: %s", err))
		return
	}
	ctx.Logger().Info(fmt.Sprintf("scheduler: found %d birthdays", len(birthdays)))
	for _, birthday := range birthdays {
		msg := tg.NewMessage(birthday.UserID, fmt.Sprintf("Birthday of %s is %s\n%s", birthday.Name, birthday.Date, birthday.Additional))
		_, err := ctx.Bot().Send(msg)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("scheduler: error: %s", err))
		}
	}
	ctx.Logger().Info("scheduler: finish job")
}
