package repo

import (
	"fmt"
	"hw_4_1/internal/entity"
	"io"
)

type UrlRepo interface {
	GetUrls() ([]string, error)
}

type ScrapeResultsRepo interface {
	SaveResults([]entity.ScrapeResult) error
}

type PageDownloader interface {
	DownloadPage(url string) (io.Reader, io.Closer, error)
}

type NotSuccessResponseCodeError struct {
	code int
}

func (e NotSuccessResponseCodeError) Error() string {
	return fmt.Sprintf("Status code: %d", e.code)
}

func (e NotSuccessResponseCodeError) StatusCode() int {
	return e.code
}
