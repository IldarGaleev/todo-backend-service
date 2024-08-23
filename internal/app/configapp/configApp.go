// Package configapp implements application configuration
package configapp

import (
	"errors"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type AppConfig struct {
	EnvMode string `yaml:"env-mode" env:"ENV_MODE" env-default:"prod"`
	Port    int    `yaml:"port" env:"PORT" env-default:"9090"`
	Dsn     string `yaml:"dsn" env:"DSN" env-require:"true"`

	SecretKey     []byte        `yaml:"secret-key" env:"SECRET_KEY" env-require:"true"`
	SecretsMaxAge time.Duration `yaml:"secrets-max-age" env:"SECRETS_MAX_AGE" env-default:"24h"`
}

// MustLoadConfig returns app configuration. Panic if failed
func MustLoadConfig(confPath string) *AppConfig {
	var appConf AppConfig

	err := cleanenv.ReadConfig(confPath, &appConf)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	return &appConf
}
