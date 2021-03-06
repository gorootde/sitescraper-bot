package main

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/PuerkitoBio/goquery"
)

type SearchResult struct {
	scraper Scraper
	url     string
	results map[string]string
}

type Scraper struct {
	config tomlSectionScraper
}

func NewScraper(config tomlSectionScraper) *Scraper {
	return &Scraper{config}

}

func (s Scraper) Scrape(resultChannel chan *SearchResult) {
	doc := s.getSource(s.config.URL)

	if s.config.Deeplinkselector == "" {
		logrus.Debug("No deeplink selector set. Switching to single page mode")
		result := s.find(doc)
		resultChannel <- &SearchResult{scraper: s, url: s.config.URL, results: result}
	} else {
		logrus.Debugf("Searching for deeplinks using selector '%s'", s.config.Deeplinkselector)
		doc.Find(s.config.Deeplinkselector).Each(func(i int, selection *goquery.Selection) {
			band, ok := selection.Attr("href")
			if ok {

				if strings.HasPrefix(band, "/") {
					url, _ := url.Parse(s.config.URL)
					band = fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, band)
				}
				logrus.Debugf("Scraping deeplink %s", band)
				src := s.getSource(band)
				result := s.find(src)
				resultChannel <- &SearchResult{scraper: s, url: band, results: result}
			}
		})
	}
}

func (s Scraper) find(src *goquery.Document) map[string]string {
	result := make(map[string]string)
	for fieldname, definition := range s.config.Fields {
		res := src.Find(definition.Selector)
		var value string
		if definition.Attribute == "" {
			value = res.Text()
		} else {
			value, _ = res.Attr(definition.Attribute)
		}

		result[fieldname] = value
	}
	return result
}

func (s Scraper) getSource(url string) *goquery.Document {
	// Request the HTML page.
	res, err := http.Get(url)
	if err != nil {
		logrus.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		logrus.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		logrus.Fatal(err)
	}
	return doc
}
