package config

import (
	"errors"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ParallelRequestsCount int    `env:"PARALLELS" env-default:"2"`
	InputFile             string `env:"INPUT" env-default:""`
	OutputFile            string `env:"OUTPUT" env-default:""`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("Не удалось прочитать параметры конфига: %w", err)
	}

	if cfg.InputFile == "" || cfg.OutputFile == "" {
		return nil, errors.New("значения env-переменных INPUT и OUTPUT должны быть указаны")
	}

	return &cfg, nil
}
