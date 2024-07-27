/*
Implements application configuration
*/
package configApp

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	EnvMode string `env:"ENV_MODE" env-default:"prod"`
	Port    int    `env:"PORT" env-default:"9090"`
}

// Returns app configuration. Panic if failed
func MustLoadConfig() *AppConfig {
	var appConf AppConfig
	err := cleanenv.ReadEnv(&appConf)
	if err != nil {
		panic(err)
	}
	return &appConf
}
