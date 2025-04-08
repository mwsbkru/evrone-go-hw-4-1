package repo

import "hw_4_1/internal/entity"

type UrlRepo interface {
	GetUrls() ([]string, error)
}

type ScrapeResultsRepo interface {
	SaveResults([]entity.ScrapeResult) error
}
