package main

import (
	"fmt"
	"log/slog"
	"os"
	"telegram-bot/internal/app"
	"telegram-bot/internal/app/context"
	"telegram-bot/internal/config"
	"telegram-bot/internal/storage/sqlite"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	cfg := config.SetupConfig()

	logger := setupLogger(cfg.Env)

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error(fmt.Sprintf("main: %v", err))
		os.Exit(1)
	}

	bot, err := tg.NewBotAPI(cfg.Telegram.BotToken)
	if err != nil {
		logger.Error(fmt.Sprintf("main: %v", err))
		os.Exit(1)
	}

	if cfg.Env == envLocal {
		bot.Debug = true
	}

	//me, err := bot.GetChat(tg.ChatInfoConfig{ChatConfig: tg.ChatConfig{ChatID: 901977201}})
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Println(me)

	updateConfig := tg.NewUpdate(0)
	updateConfig.Timeout = cfg.Telegram.RequestTimeout
	updates := bot.GetUpdatesChan(updateConfig)

	ctx := context.New(bot, storage, logger)

	app.GetUpdates(updates, ctx)
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

	}
	return logger
}
