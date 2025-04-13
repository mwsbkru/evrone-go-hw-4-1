package repo

import (
	"io"
	"strings"
	"time"
)

var pages map[string]string = map[string]string{
	"https://httpd.apache.org/docs/2.2/en/content-negotiation.html#algorithm":      "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Apache DESC\"><title>Apache TITLE</title></head><body id=\"manual-page\"><h1>Apache</h1></body></html>",
	"https://developer.mozilla.org/en-US/docs/Web/HTTP/Guides/Content_negotiation": "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Mozilla DESC\"><title>Mozilla TITLE</title></head><body id=\"manual-page\"><h1>Mozilla</h1></body></html>",
	"https://apple.com/Content_negotiation":                                        "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Apple DESC\"><title>Apple TITLE</title></head><body id=\"manual-page\"><h1>Apple</h1></body></html>",
	"https://za.com/Content_negotiation":                                           "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Za DESC\"><title>Za TITLE</title></head><body id=\"manual-page\"><h1>Za</h1></body></html>",
	"https://qwe.com/Content_negotiation":                                          "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"qwe DESC\"><title>Qwe TITLE</title></head><body id=\"manual-page\"><h1>Qwe</h1></body></html>",
}

type DummyPageDownloader struct {
}

func NewDummyPageDownloader() *DummyPageDownloader {
	return &DummyPageDownloader{}
}

func (d *DummyPageDownloader) DownloadPage(url string) (io.Reader, error) {
	time.Sleep(1*time.Second + time.Duration(len(url)*50))
	if page, ok := pages[url]; ok {
		return strings.NewReader(page), nil
	}

	return nil, NotSuccessResponseCodeError{code: 404}
}
