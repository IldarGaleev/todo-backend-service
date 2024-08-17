// Package configapp implements application configuration
package configapp

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	EnvMode string `env:"ENV_MODE" env-default:"prod"`
	Port    int    `env:"PORT" env-default:"9090"`
	Dsn     string `env:"DSN" env-require:"true"`

	SecretKey     []byte        `env:"SECRET_KEY" env-require:"true"`
	SecretsMaxAge time.Duration `env:"SECRETS_MAX_AGE" env-default:"24h"`
}

// MustLoadConfig returns app configuration. Panic if failed
func MustLoadConfig() *AppConfig {
	//TODO: add loading config from file
	var appConf AppConfig
	err := cleanenv.ReadEnv(&appConf)
	if err != nil {
		panic(err)
	}
	return &appConf
}
