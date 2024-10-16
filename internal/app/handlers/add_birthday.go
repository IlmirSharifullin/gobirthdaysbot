package handlers

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/states"
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
		ctx.UpdateData(user.ID, "name", name)

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
	ctx.UpdateData(user.ID, "dateOfBirth", dateOfBirth)

	msg := tg.NewMessage(user.ID, "Now enter an additional information about the person:")
	ctx.SetState(user.ID, states.BirthdayAdditionalState)
	_, err = ctx.Bot().Send(msg)
	return err
}

func BirthdayAdditional(ctx common.Context, u tg.Update) error {
	return nil
}
