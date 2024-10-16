package context

import (
	"fmt"
	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"telegram-bot/internal/app/common"
	"telegram-bot/internal/app/states"
	"telegram-bot/internal/kvstorage"
	mapped "telegram-bot/internal/kvstorage/map"
	"telegram-bot/internal/storage"
)

type Context struct {
	bot       *tg.BotAPI
	db        storage.Storage
	kvstorage kvstorage.KVStorage
	logger    *slog.Logger
	State     map[int64]*common.State
}

func New(bot *tg.BotAPI, db storage.Storage, logger *slog.Logger) *Context {
	return &Context{bot: bot, db: db, kvstorage: mapped.New(), logger: logger, State: map[int64]*common.State{}}
}

func (ctx *Context) Bot() *tg.BotAPI {
	return ctx.bot
}

func (ctx *Context) Db() storage.Storage {
	return ctx.db
}

func (ctx *Context) KVStorage() kvstorage.KVStorage {
	return ctx.kvstorage
}

func (ctx *Context) Logger() *slog.Logger {
	return ctx.logger
}

func (ctx *Context) Clear(ID int64) {
	ctx.kvstorage.Clear(ID)
}

func (ctx *Context) UpdateData(ID int64, dataKey string, dataValue any) {
	ctx.kvstorage.Update(ID, dataKey, dataValue)
}

func (ctx *Context) GetData(ID int64, dataKey string) any {
	return ctx.kvstorage.Get(ID, dataKey)
}

func (ctx *Context) SetState(ID int64, s *common.State) {
	ctx.State[ID] = s
}

func (ctx *Context) Serve(u tg.Update) error {
	user := u.SentFrom()
	if user == nil {
		return common.ErrNoUser
	}
	state, ok := ctx.State[user.ID]
	if !ok {
		ctx.State[user.ID] = states.AnyState
		state = states.AnyState
	}

	if state.Filter(u) {
		err := state.Handler(ctx, u)
		ctx.logger.Error(fmt.Sprintf("%v", err))
		return ctx.ErrorHandler(u, err)
	} else {
		err := state.ElseFunc(ctx, u)
		ctx.logger.Error(fmt.Sprintf("%v", err))
		return ctx.ErrorHandler(u, err)
	}
}

func (ctx *Context) ErrorHandler(u tg.Update, err error) error {
	if err == nil {
		return nil
	}
	user := u.SentFrom()
	if user == nil {
		return nil
	}
	ctx.SetState(user.ID, states.AnyState)
	msg := tg.NewMessage(user.ID, "Sorry, an error occurred.")
	_, err = ctx.bot.Send(msg)
	return err
}
