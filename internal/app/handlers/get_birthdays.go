package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
)

func GetAllBirthdays(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}

	birthdays, err := ctx.Db().GetBirthdays(user.ID)
	if err != nil {
		return err
	}
	ctx.Logger().Info(fmt.Sprintf("found %d birthdays", len(birthdays)))

	for _, birthday := range birthdays {
		msg := common.GetBirthdayCard(birthday)
		_, err := ctx.Bot().Send(msg)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("error: %s", err))
		}
	}
	return nil
}

func GetNextBirthday(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}

	birthdays, err := ctx.Db().GetNextBirthdays(user.ID)
	if err != nil {
		return err
	}
	for _, birthday := range birthdays {
		msg := common.GetBirthdayCard(birthday)
		if _, err = ctx.Bot().Send(msg); err != nil {
			ctx.Logger().Error(fmt.Sprintf("error: %s", err))
		}
	}
	return nil
}
