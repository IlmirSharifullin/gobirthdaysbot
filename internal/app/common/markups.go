package common

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"telegram-bot/internal/storage"
	"telegram-bot/pkg/callback_data"
)

var KeyboardMenu = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Add a birthday"),
		tg.NewKeyboardButton("Get all birthdays"),
		tg.NewKeyboardButton("Get next birthday"),
	),
)

func MakeCardMarkup(birthday *storage.Birthday) tg.InlineKeyboardMarkup {
	deleteCD := callback_data.CallbackData{Prefix: "delete_", Sep: "_", KVSep: ":", Data: map[string]string{}}
	deleteCD.Data["id"] = strconv.FormatInt(birthday.ID, 10)

	kb := tg.NewInlineKeyboardMarkup(
		tg.NewInlineKeyboardRow(
			tg.NewInlineKeyboardButtonData("Delete 🗑", deleteCD.String()),
		),
	)
	return kb
}
