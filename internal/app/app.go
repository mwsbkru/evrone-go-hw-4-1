package app

import (
	"context"
	"hw_4_1/config"
	"hw_4_1/internal/repo"
	"hw_4_1/internal/usecase"
	"log/slog"
)

func Run(ctx context.Context, cfg *config.Config) {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	urlRepo := repo.NewFileUrlRepo(cfg.InputFile)
	resultsRepo := repo.NewCsvScrapeResultRepo(cfg.OutputFile)
	downloader := repo.NewSimplePageDownloader()
	scraper := usecase.NewScraper(*cfg, urlRepo, resultsRepo, downloader)

	err := scraper.Scrape(ctx)
	if err != nil {
		slog.Error("Не удалось выполнить сканирование ссылок", slog.String("error", err.Error()))
	}
}
