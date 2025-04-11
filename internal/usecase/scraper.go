package usecase

import (
	"fmt"
	"hw_4_1/internal/entity"
	"hw_4_1/internal/repo"
	"hw_4_1/internal/service"
)

type Scraper struct {
	urlRepo     repo.UrlRepo
	resultsRepo repo.ScrapeResultsRepo
	downloader  repo.PageDownloader
}

func NewScraper(urlRepo repo.UrlRepo, resultsRepo repo.ScrapeResultsRepo, downloader repo.PageDownloader) Scraper {
	return Scraper{urlRepo: urlRepo, resultsRepo: resultsRepo, downloader: downloader}
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
	results := make([]entity.ScrapeResult, 0, len(urls))
	urlParser := service.NewHtmlParser()

	for _, url := range urls {
		body, err := s.downloader.DownloadPage(url)
		if err != nil {
			return results, fmt.Errorf("не удалось просканировать переданные url: %w", err)
		}

		result, err := urlParser.ParseHtml(body, url)
		if err != nil {
			return results, fmt.Errorf("не удалось распарсить тело страницы %s: %w", url, err)
		}

		results = append(results, result)
	}
	return results, nil
}
