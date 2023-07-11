package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	BaseYandexDiskAPIUrl string `env:"BASE_YANDEX_DISK_API_URL" env-required:"true"`
	TG_Bot_Token         string `env:"TG_BOT_TOKEN" env-required:"true"`
}

func MustLoad(configPath string) *Config {

	// check if file exist
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file is not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
