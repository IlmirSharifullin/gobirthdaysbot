package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/storage"
)

func Echo(u tg.Update, bot *tg.BotAPI, db storage.Storage) error {
	msg := tg.NewMessage(u.Message.From.ID, u.Message.Text)
	if _, err := bot.Send(msg); err != nil {
		return err
	}
	return nil
}
