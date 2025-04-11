package app

import (
	"hw_4_1/config"
	"hw_4_1/internal/repo"
	"hw_4_1/internal/usecase"
	"log/slog"
)

func Run(cfg *config.Config) {
	urlRepo := repo.NewFileUrlRepo(cfg.InputFile)
	resultsRepo := repo.NewCsvScrapeResultRepo(cfg.OutputFile)
	downloader := repo.NewDummyPageDownloader()
	scraper := usecase.NewScraper(urlRepo, resultsRepo, downloader)

	err := scraper.Scrape()
	if err != nil {
		slog.Error("Не удалось выполнить сканирование ссылок", slog.String("error", err.Error()))
	}
}
