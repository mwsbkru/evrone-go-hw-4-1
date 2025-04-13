package usecase

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/semaphore"
	"hw_4_1/config"
	"hw_4_1/internal/entity"
	"hw_4_1/internal/repo"
	"hw_4_1/internal/service"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Scraper struct {
	cfg         config.Config
	urlRepo     repo.UrlRepo
	resultsRepo repo.ScrapeResultsRepo
	downloader  repo.PageDownloader
	htmlParser  *service.HtmlParser
}

func NewScraper(cfg config.Config, urlRepo repo.UrlRepo, resultsRepo repo.ScrapeResultsRepo, downloader repo.PageDownloader) Scraper {
	return Scraper{cfg: cfg, urlRepo: urlRepo, resultsRepo: resultsRepo, downloader: downloader, htmlParser: service.NewHtmlParser()}
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
	ctx := context.TODO()
	var wg sync.WaitGroup
	// TODO воткнцть консткст
	sem := semaphore.NewWeighted(int64(s.cfg.ParallelRequestsCount))
	results := make([]entity.ScrapeResult, 0, len(urls))
	resultChan := make(chan entity.ScrapeResult)

	for _, url := range urls {
		url := url
		wg.Add(1)
		go s.scrapeUrl(ctx, url, &wg, sem, resultChan)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		results = append(results, result)
	}
	return results, nil
}

func (s *Scraper) scrapeUrl(ctx context.Context, url string, wg *sync.WaitGroup, sem *semaphore.Weighted, resultChan chan<- entity.ScrapeResult) {
	slog.Debug("1. Начало сканирования страницы", slog.String("url", url))
	result := entity.ScrapeResult{Date: time.Now(), Url: url}
	defer func() {
		resultChan <- result
		wg.Done()
	}()

	err := sem.Acquire(ctx, 1)
	if err != nil {
		slog.Error("Не удалось запустить параллельную обработку запросов", slog.String("url", url), slog.String("error", err.Error()))
		return
	}

	body, err := s.downloader.DownloadPage(url)
	sem.Release(1)
	var notSuccessResponseCodeErr repo.NotSuccessResponseCodeError
	if err != nil {
		if errors.As(err, &notSuccessResponseCodeErr) {
			result.StatusCode = notSuccessResponseCodeErr.StatusCode()
			slog.Error(notSuccessResponseCodeErr.Error(), slog.String("url", url))
		} else {
			slog.Error("не удалось просканировать переданный url", slog.String("error", err.Error()), slog.String("url", url))
		}
		return
	}
	result.StatusCode = http.StatusOK

	slog.Debug("2. Начало парсинга страницы", slog.String("url", url))
	pageData, err := s.htmlParser.ParseHtml(body)
	if err != nil {
		slog.Error("не удалось распарсить тело страницы", slog.String("error", err.Error()), slog.String("url", url))
		return
	}

	slog.Debug("3. Сканирование страницы завершено", slog.String("url", url))
	result.SuccessAttempt += 1
	result.Title = pageData.Title
	result.Description = pageData.Description
}
