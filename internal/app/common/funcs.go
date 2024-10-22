package common

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/storage"
	"telegram-bot/pkg/dates"
)

var MsgNotNil = func(u tg.Update) bool {
	return u.Message != nil
}

func GetBirthdayCard(birthday *storage.Birthday) string {
	years, days := dates.CalculateDate(birthday.Date)
	text := fmt.Sprintf("%s\n\nBirthday is %s (%s turns %d %s)\n%s", birthday.Name, birthday.Date.Format("02.01.2006"), birthday.Name, years, days, birthday.Additional)
	return text
}
