package handlers

import (
	"errors"
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/states"
	"telegram-bot/internal/storage"
)

func StartCommand(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	ctx.SetState(user.ID, states.AnyState)

	if user == nil {
		return nil
	}
	dbUser, err := ctx.Db().GetUser(user.ID)
	if err != nil {
		if !errors.Is(err, storage.ErrUserNotExists) {
			return err
		}
		err = ctx.Db().InsertUser(user.ID, user.UserName)
		if err != nil {
			return err
		}
		msg := tg.NewMessage(user.ID, "Hello to my new bot developed on Golang!")
		msg.ReplyMarkup = common.KeyboardMenu
		_, err = ctx.Bot().Send(msg)
		return err
	} else {
		msg := tg.NewMessage(user.ID, fmt.Sprintf("Hello again, %s!", getUsername(dbUser.Username)))
		msg.ReplyMarkup = common.KeyboardMenu
		_, err = ctx.Bot().Send(msg)
		return err
	}
}

func getUsername(s string) string {
	if s == "" {
		return "Guest"
	}
	return s
}
