package states

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
)

var AnyState = &common.State{Filter: func(u tg.Update) bool { return true }}

var BirthdayNameState = &common.State{Filter: common.MsgNotNil}
var BirthdayDateState = &common.State{Filter: common.MsgNotNil}
var BirthdayAdditionalState = &common.State{Filter: common.MsgNotNil}
