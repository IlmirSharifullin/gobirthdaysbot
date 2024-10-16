package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string         `env:"ENV" env-default:"local"`
	StoragePath string         `env:"STORAGE_PATH" end-default:"db/storage.db"`
	Telegram    TelegramConfig `env-prefix:"TELEGRAM_"`
}

type TelegramConfig struct {
	BotToken       string `env:"BOT_TOKEN" required:"true"`
	RequestTimeout int    `env:"REQUEST_TIMEOUT" env-default:"30"`
}

func SetupConfig() *Config {
	cfg := Config{}
	_ = cleanenv.ReadConfig(".env", &cfg)
	return &cfg
}
