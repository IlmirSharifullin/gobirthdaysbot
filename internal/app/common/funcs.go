package common

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

var MsgNotNil = func(u tg.Update) bool {
	return u.Message != nil
}
