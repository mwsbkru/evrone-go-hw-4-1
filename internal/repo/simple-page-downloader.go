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

func (d *SimplePageDownloader) DownloadPage(url string) (io.Reader, io.Closer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, fmt.Errorf("не удалось выполнить http-запрос: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, &NotSuccessResponseCodeError{code: resp.StatusCode}
	}

	return resp.Body, resp.Body, nil
}
