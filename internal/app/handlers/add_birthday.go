package handlers

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/states"
	"telegram-bot/internal/storage"
	"time"
)

func BirthdayAddCommand(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}

	msg := tg.NewMessage(user.ID, "Enter the name of the person whose birthday is coming up")
	_, err := ctx.Bot().Send(msg)

	ctx.SetState(user.ID, states.BirthdayNameState)
	ctx.UpdateData(user.ID, "birthday", &storage.Birthday{})

	return err
}

func BirthdayName(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}
	defer ctx.SetState(user.ID, states.BirthdayDateState)
	name := u.Message.Text
	if name == "" {
		msg := tg.NewMessage(user.ID, "Enter the name of the person whose birthday is coming up")
		_, err := ctx.Bot().Send(msg)
		return err
	} else {
		birthday := ctx.GetData(user.ID, "birthday")
		switch birthday.(type) {
		case *storage.Birthday:
			birthday.(*storage.Birthday).Name = name
		default:
			return common.ErrBirthdayTypeError
		}

		msg := tg.NewMessage(user.ID, "Enter this person's date of birth (format: DD.MM.YYYY)")
		_, err := ctx.Bot().Send(msg)
		return err
	}
}

func BirthdayDate(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}

	dateOfBirth, err := time.Parse(time.DateOnly, u.Message.Text)
	if err != nil {
		return common.ErrNotDate
	}

	birthday := ctx.GetData(user.ID, "birthday")
	switch birthday.(type) {
	case *storage.Birthday:
		birthday.(*storage.Birthday).Date = dateOfBirth
	default:
		return common.ErrBirthdayTypeError
	}

	msg := tg.NewMessage(user.ID, "Now enter an additional information about the person:")
	ctx.SetState(user.ID, states.BirthdayAdditionalState)
	_, err = ctx.Bot().Send(msg)
	return err
}

func BirthdayAdditional(ctx common.Context, u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}
	defer ctx.SetState(user.ID, states.AnyState)

	text := u.Message.Text
	birthday := ctx.GetData(user.ID, "birthday")
	switch birthday.(type) {
	case *storage.Birthday:
		birthday.(*storage.Birthday).Additional = text
	default:
		return common.ErrBirthdayTypeError
	}

	_birthday := birthday.(*storage.Birthday)
	_birthday.UserID = user.ID

	err := ctx.Db().InsertBirthday(_birthday)
	ctx.Logger().Info(fmt.Sprintf("%v", _birthday))
	if err != nil {
		return err
	}
	msg := tg.NewMessage(user.ID, "Birthday successfully added!")
	_, err = ctx.Bot().Send(msg)
	return err
}
