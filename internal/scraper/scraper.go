package scraper

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Metadata struct {
	Title       string
	Description string
	Image       string
}

func FetchMetadata(url string) (*Metadata, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	title := strings.TrimSpace(doc.Find("title").Text())
	desc, _ := doc.Find("meta[name='description']").Attr("content")
	image, _ := doc.Find("meta[property='og:image']").Attr("content")

	return &Metadata{
		Title:       title,
		Description: desc,
		Image:       image,
	}, nil
}
