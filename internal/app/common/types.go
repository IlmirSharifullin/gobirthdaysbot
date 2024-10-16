package common

import (
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"telegram-bot/internal/kvstorage"
	"telegram-bot/internal/storage"
)

type Handler func(ctx Context, u tg.Update) error

type State struct {
	Filter   func(u tg.Update) bool
	Handler  Handler
	ElseFunc Handler
}

type Context interface {
	Bot() *tg.BotAPI
	Db() storage.Storage
	KVStorage() kvstorage.KVStorage
	Logger() *slog.Logger
	Clear(ID int64)
	SetState(ID int64, s *State)
	UpdateData(ID int64, dataKey string, dataValue any)
	Serve(u tg.Update) error
	ErrorHandler(u tg.Update, err error) error
}
