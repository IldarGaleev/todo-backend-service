package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type appConfig struct {
	ConfigStr string `env:"CONFIG_STR" env-default:"Default str"`
}

var AppConfig appConfig

func LoadConfig() error {
	err := cleanenv.ReadEnv(&AppConfig)
	if err != nil {
		return fmt.Errorf("Load config error")
	}
	return nil
}
