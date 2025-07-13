package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"go.uber.org/zap"
)

type Config struct {
	Port  int `env:"PORT" env-default:"8080"`
	DBCfg DBConfig
}

type DBConfig struct {
	Host     string `env:"DB_HOST"`
	Port     string `env:"DB_PORT"`
	User     string `env:"DB_USER"`
	Password string `env:"DB_PASSWORD"`
	Name     string `env:"DB_NAME"`
	SSLMode  string `env:"DB_SSLMODE" env-default:"disable"`
}

func InitConfig(logger *zap.Logger) *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		logger.Fatal("Failed to initialize config: ", zap.Error(err))
	}

	return &cfg
}
