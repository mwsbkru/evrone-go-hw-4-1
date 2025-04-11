package repo

import (
	"fmt"
	"io"
	"net/http"
)

type SimplePageDownloader struct {
}

func NewSimplePageDownloader() *SimplePageDownloader {
	return &SimplePageDownloader{}
}

func (d *SimplePageDownloader) DownloadPage(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("не удалось выполнить http-запрос: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, NotSuccessResponseCodeError{code: resp.StatusCode}
	}

	return resp.Body, nil
}
