package service

import (
	"errors"
	"hw_4_1/internal/entity"
	"io"
	"strings"
	"testing"
)

type readerWithError struct {
}

func (r *readerWithError) Read(_ []byte) (int, error) {
	return 0, errors.New("reader with error")
}

func TestHtmlParser_ParseHtml(t *testing.T) {
	parser := NewHtmlParser()

	testCases := []struct {
		name             string
		input            io.Reader
		expectedPageData entity.PageData
		wantError        bool
	}{
		{
			name:  "Normal body",
			input: strings.NewReader("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"><html xmlns=\"http://www.w3.org/1999/xhtml\" lang=\"en\" xml:lang=\"en\"><head><meta name=\"description\" content=\"qwe DESC\"><title>Qwe TITLE</title></head><body id=\"manual-page\"><h1>Qwe</h1></body></html>"),
			expectedPageData: entity.PageData{
				Title:       "Qwe TITLE",
				Description: "qwe DESC",
			},
		},
		{
			name:  "Empty body",
			input: strings.NewReader(""),
			expectedPageData: entity.PageData{
				Title:       "",
				Description: "",
			},
		},
		{
			name:  "Broken html",
			input: strings.NewReader("qwerty"),
			expectedPageData: entity.PageData{
				Title:       "",
				Description: "",
			},
		},
		{
			name:      "Reader error",
			input:     &readerWithError{},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		actualResult, err := parser.ParseHtml(tc.input)

		if (err != nil) != tc.wantError {
			t.Errorf("Test case - %s: want error", tc.name)
		}

		if actualResult != tc.expectedPageData {
			t.Errorf("Test case - %s: want title - %v; got - %v", tc.name, tc.expectedPageData, actualResult)
		}
	}
}
