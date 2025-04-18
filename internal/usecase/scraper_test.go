package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"hw_4_1/config"
	"hw_4_1/internal/entity"
	"hw_4_1/internal/repo"
	"log/slog"
	"net/http"
	"strings"
	"testing"
)

type mocksCollection struct {
	urlRepo     *repo.MockUrlRepo
	downloader  *repo.MockPageDownloader
	resultsRepo *stubScrapeResultsRepo
}

type stubScrapeResultsRepo struct {
	results     []entity.ScrapeResult
	errToReturn error
}

func (s *stubScrapeResultsRepo) SaveResults(results []entity.ScrapeResult) error {
	s.results = results
	return s.errToReturn
}

func TestScraper_Scrape(t *testing.T) {
	slog.SetLogLoggerLevel(10)

	cfg := config.Config{
		ParallelRequestsCount: 2,
		RetryTimeoutSeconds:   0,
		MaxAttemptCount:       3,
	}

	tests := []struct {
		name             string
		setup            func(*mocksCollection)
		checkResults     func(*stubScrapeResultsRepo, *testing.T, string)
		wantErr          bool
		expectedPageData entity.ScrapeResult
	}{
		{
			name:    "success scrape",
			wantErr: false,
			setup: func(mocksCollection *mocksCollection) {
				pageBody := "<html><head><meta name=\"description\" content=\"qwe DESC\"><title>Qwe TITLE</title></head><body id=\"manual-page\"><h1>Qwe</h1></body></html>"
				mocksCollection.urlRepo.EXPECT().GetUrls().Return([]string{"http://url.test"}, nil)
				mocksCollection.downloader.EXPECT().DownloadPage("http://url.test").Return(strings.NewReader(pageBody), nil, nil)
			},
			checkResults: func(resultsRepo *stubScrapeResultsRepo, t *testing.T, testName string) {
				if len(resultsRepo.results) != 1 {
					t.Errorf("Test case - %s: wrong results count. Want 1, got - %d", testName, len(resultsRepo.results))
				}

				result := resultsRepo.results[0]

				if result.Url != "http://url.test" {
					t.Errorf("Test case - %s: wrong result url. Want %s, got - %s", testName, "http://url.test", result.Url)
				}

				if result.StatusCode != http.StatusOK {
					t.Errorf("Test case - %s: wrong result status code. Want %d, got - %d", testName, http.StatusOK, result.StatusCode)
				}

				if result.SuccessAttempt != 1 {
					t.Errorf("Test case - %s: wrong result success attempt. Want %d, got - %d", testName, 1, result.SuccessAttempt)
				}
			},
		},
		{
			name:    "success scrape with retry",
			wantErr: false,
			setup: func(mocksCollection *mocksCollection) {
				pageBody := "<html><head><meta name=\"description\" content=\"qwe DESC\"><title>Qwe TITLE</title></head><body id=\"manual-page\"><h1>Qwe</h1></body></html>"
				mocksCollection.urlRepo.EXPECT().GetUrls().Return([]string{"http://url.test"}, nil)
				mocksCollection.downloader.EXPECT().DownloadPage("http://url.test").Return(strings.NewReader(pageBody), nil, &repo.NotSuccessResponseCodeError{Code: 404})
				mocksCollection.downloader.EXPECT().DownloadPage("http://url.test").Return(strings.NewReader(pageBody), nil, errors.New("oops"))
				mocksCollection.downloader.EXPECT().DownloadPage("http://url.test").Return(strings.NewReader(pageBody), nil, nil)
			},
			checkResults: func(resultsRepo *stubScrapeResultsRepo, t *testing.T, testName string) {
				if len(resultsRepo.results) != 1 {
					t.Errorf("Test case - %s: wrong results count. Want 1, got - %d", testName, len(resultsRepo.results))
				}

				result := resultsRepo.results[0]

				if result.StatusCode != http.StatusOK {
					t.Errorf("Test case - %s: wrong result status code. Want %d, got - %d", testName, http.StatusOK, result.StatusCode)
				}

				if result.SuccessAttempt != 3 {
					t.Errorf("Test case - %s: wrong result success attempt. Want %d, got - %d", testName, 3, result.SuccessAttempt)
				}
			},
		},
		{
			name:    "scrape with bad response",
			wantErr: false,
			setup: func(mocksCollection *mocksCollection) {
				mocksCollection.urlRepo.EXPECT().GetUrls().Return([]string{"http://url.test"}, nil)
				mocksCollection.downloader.EXPECT().DownloadPage("http://url.test").Times(3).Return(nil, nil, errors.New("oops"))

			},
			checkResults: func(resultsRepo *stubScrapeResultsRepo, t *testing.T, testName string) {
				if len(resultsRepo.results) != 1 {
					t.Errorf("Test case - %s: wrong results count. Want 1, got - %d", testName, len(resultsRepo.results))
				}

				result := resultsRepo.results[0]

				if result.StatusCode != 0 {
					t.Errorf("Test case - %s: wrong result status code. Want %d, got - %d", testName, 0, result.StatusCode)
				}

				if result.SuccessAttempt != 0 {
					t.Errorf("Test case - %s: wrong result success attempt. Want %d, got - %d", testName, 0, result.SuccessAttempt)
				}
			},
		},
		{
			name:    "scrape with failed get urls",
			wantErr: true,
			setup: func(mocksCollection *mocksCollection) {
				mocksCollection.urlRepo.EXPECT().GetUrls().Return([]string{}, errors.New("oops"))
			},
			checkResults: nil,
		},
		{
			name:    "scrape with failed save results",
			wantErr: true,
			setup: func(mocksCollection *mocksCollection) {
				mocksCollection.urlRepo.EXPECT().GetUrls().Return([]string{}, nil)
				mocksCollection.resultsRepo.errToReturn = errors.New("oops")
			},
			checkResults: nil,
		},
	}

	for _, tc := range tests {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mocks := mocksCollection{
			urlRepo:     repo.NewMockUrlRepo(ctrl),
			downloader:  repo.NewMockPageDownloader(ctrl),
			resultsRepo: &stubScrapeResultsRepo{},
		}

		if tc.setup != nil {
			tc.setup(&mocks)
		}

		scraper := NewScraper(cfg, mocks.urlRepo, mocks.resultsRepo, mocks.downloader)

		if err := scraper.Scrape(context.Background()); (err != nil) != tc.wantErr {
			t.Errorf("Test case - %s: want error - %v", tc.name, tc.wantErr)
		}

		if tc.checkResults != nil {
			tc.checkResults(mocks.resultsRepo, t, tc.name)
		}
	}
}
