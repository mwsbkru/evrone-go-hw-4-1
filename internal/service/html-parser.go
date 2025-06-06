package service

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"hw_4_1/internal/entity"
	"io"
)

type HtmlParser struct {
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

func (p *HtmlParser) ParseHtml(html io.Reader) (entity.PageData, error) {
	doc, err := goquery.NewDocumentFromReader(html)
	if err != nil {
		return entity.PageData{}, fmt.Errorf("Не удалось распарсить url: %w", err)
	}

	title := doc.Find("title").Text()
	description := doc.Find("meta[name=\"description\"]").AttrOr("content", "")

	return entity.PageData{
		Title:       title,
		Description: description,
	}, nil
}
