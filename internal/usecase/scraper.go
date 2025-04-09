package usecase

import (
	"fmt"
	"hw_4_1/internal/entity"
	"hw_4_1/internal/repo"
	"time"
)

type Scraper struct {
	urlRepo     repo.UrlRepo
	resultsRepo repo.ScrapeResultsRepo
}

func NewScraper(urlRepo repo.UrlRepo, resultsRepo repo.ScrapeResultsRepo) Scraper {
	return Scraper{urlRepo: urlRepo, resultsRepo: resultsRepo}
}

func (s *Scraper) Scrape() error {
	urls, err := s.urlRepo.GetUrls()
	if err != nil {
		return fmt.Errorf("Не удалось прочитать список ссылок для сканирования: %w", err)
	}

	scrapeResults, err := s.processScrape(urls)
	if err != nil {
		return fmt.Errorf("Не удалось просканировать ссылки: %w", err)
	}

	err = s.resultsRepo.SaveResults(scrapeResults)
	if err != nil {
		return fmt.Errorf("Не удалось сохранить результаты сканирования: %w", err)
	}

	return nil
}

func (s *Scraper) processScrape(urls []string) ([]entity.ScrapeResult, error) {
	result := make([]entity.ScrapeResult, 0, len(urls))

	for _, url := range urls {
		scrapeResult := entity.ScrapeResult{
			Date:           time.Now(),
			Url:            url,
			StatusCode:     len(url),
			Title:          "title",
			Description:    "descr",
			SuccessAttempt: 0,
		}
		result = append(result, scrapeResult)
	}
	return result, nil
}
