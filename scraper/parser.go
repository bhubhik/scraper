package scraper

import (
	"errors"
	"io"

	"github.com/PuerkitoBio/goquery"
)

func ParseTop10(body io.ReadCloser) ([]string, error) {
	if body == nil {
		return nil, errors.New("Cannot parse empty body")
	}

	doc, err := goquery.NewDocumentFromReader(body)
	if err != nil {
		return nil, err
	}

	items := []string{}

	doc.Find("ul li h3").Each(func(i int, s *goquery.Selection) {
		if i >= 10 {
			return
		}
		title := s.Text()
		items = append(items, title)
	})
	return items, nil
}
