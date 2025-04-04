package main

import (
	"hw_4_1/config"
	"hw_4_1/internal/app"
	"log/slog"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Не удалось инициализировать конфиг", slog.String("error", err.Error()))
	}

	app.Run(cfg)

}
