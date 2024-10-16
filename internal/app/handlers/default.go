package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/storage"
)

func IDontUnderstand(u tg.Update, bot tg.BotAPI, db storage.Storage) error {
	user := u.SentFrom()
	if user != nil {
		_, err := bot.Send(tg.NewMessage(user.ID, "Sorry, i don`t understand you! ;("))
		return err
	}
	return nil
}
