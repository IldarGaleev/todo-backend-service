package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
)

type appConfig struct {
	Port int `env:"PORT" env-default:"9090"`
}

var AppConfig appConfig

func LoadConfig() error {
	err := cleanenv.ReadEnv(&AppConfig)
	if err != nil {
		return fmt.Errorf("Load config error")
	}
	return nil
}
