package repo

import (
	"io"
	"strings"
)

var pages map[string]string = map[string]string{
	"https://httpd.apache.org/docs/2.2/en/content-negotiation.html#algorithm":      "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Apache DESC\"><title>Apache TITLE</title></head><body id=\"manual-page\"><h1>Apache</h1></body></html>",
	"https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Content_negotiation": "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Mozilla DESC\"><title>Mozilla TITLE</title></head><body id=\"manual-page\"><h1>Mozilla</h1></body></html>",
}

type DummyPageDownloader struct {
}

func NewDummyPageDownloader() *DummyPageDownloader {
	return &DummyPageDownloader{}
}

func (d *DummyPageDownloader) DownloadPage(url string) (io.Reader, error) {
	if page, ok := pages[url]; ok {
		return strings.NewReader(page), nil
	}

	return nil, NotSuccessResponseCodeError{code: 404}
}
