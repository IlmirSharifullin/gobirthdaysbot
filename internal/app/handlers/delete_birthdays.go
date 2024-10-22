package handlers

import (
	"errors"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/storage"
	"telegram-bot/pkg/callback_data"
)

func DeleteBirthday(ctx common.Context, u tg.Update, callbackData callback_data.CallbackData) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}

	if u.CallbackData() == "" {
		return nil
	}

	var birthdayID int64
	var err error
	if v, ok := callbackData.Data["id"]; !ok {
		return common.ErrCallbackDataNoKey
	} else {
		birthdayID, err = strconv.ParseInt(v, 10, 0)
		if err != nil {
			return common.ErrCallbackDataWrongType
		}
	}

	err = ctx.Db().DeleteBirthday(birthdayID, user.ID)
	if err != nil {
		if errors.Is(err, storage.ErrNoDeleted) {
			msg := tg.NewMessage(user.ID, "This birthday does not exist or is not yours")
			_, err = ctx.Bot().Send(msg)
			return err
		}
		return err
	} else {
		callback := tg.NewCallback(u.CallbackQuery.ID, "This birthday successfully deleted")
		msg := tg.NewDeleteMessage(user.ID, u.CallbackQuery.Message.MessageID)
		if _, err := ctx.Bot().Request(callback); err != nil {
			return err
		}
		if _, err := ctx.Bot().Request(msg); err != nil {
			return err
		}
	}
	return nil
}
