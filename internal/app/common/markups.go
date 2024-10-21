package common

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var KeyboardMenu = tg.NewReplyKeyboard(
	tg.NewKeyboardButtonRow(
		tg.NewKeyboardButton("Add a birthday"),
		tg.NewKeyboardButton("Get all birthdays"),
	),
)
