package crawl

import "github.com/gocolly/colly/v2"

type BaseCrawler struct {
	Site      string
	Domains   []string
	Url       string
	Collector *colly.Collector
}
