package mycron

import (
	"fmt"
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
		msg := common.GetBirthdayCard(birthday)
		_, err := ctx.Bot().Send(msg)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("scheduler: error: %s", err))
		}
	}
	ctx.Logger().Info("scheduler: finish job")
}
