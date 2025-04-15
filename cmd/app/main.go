package main

import (
	"context"
	"hw_4_1/config"
	"hw_4_1/internal/app"
	"log/slog"
)

func main() {
	ctx := context.Background()
	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error("Не удалось инициализировать конфиг", slog.String("error", err.Error()))
	}

	app.Run(ctx, cfg)

}
