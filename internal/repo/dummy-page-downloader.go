package repo

import (
	"io"
	"strings"
	"time"
)

var pages map[string]string = map[string]string{
	"https://aq.ru/Content_negotiation":                     "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Aq DESC\"><title>Aq TITLE</title></head><body id=\"manual-page\"><h1>Aq</h1></body></html>",
	"https://sw.ru/Content_negotiation/Content_negotiation": "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"sw DESC\"><title>Sw TITLE</title></head><body id=\"manual-page\"><h1>Sw</h1></body></html>",
	"https://ap.com/Content_negotiation":                    "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Ap DESC\"><title>Ap TITLE</title></head><body id=\"manual-page\"><h1>Ap</h1></body></html>",
	"https://za.com/Content_negotiation":                    "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"Za DESC\"><title>Za TITLE</title></head><body id=\"manual-page\"><h1>Za</h1></body></html>",
	"https://qwe.com/Content_negotiation":                   "<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"qwe DESC\"><title>Qwe TITLE</title></head><body id=\"manual-page\"><h1>Qwe</h1></body></html>",
}

type DummyPageDownloader struct {
	counter int
}

func NewDummyPageDownloader() *DummyPageDownloader {
	return &DummyPageDownloader{}
}

func (d *DummyPageDownloader) DownloadPage(url string) (io.Reader, io.Closer, error) {
	time.Sleep(1*time.Second + time.Duration(len(url)*50))
	d.counter++

	if d.counter%2 != 0 {
		return nil, nil, &NotSuccessResponseCodeError{code: 404}
	}
	if page, ok := pages[url]; ok {
		return strings.NewReader(page), nil, nil
	}

	return nil, nil, &NotSuccessResponseCodeError{code: 404}
}
