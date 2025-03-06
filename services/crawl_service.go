package services

import (
	"fmt"
	"jobgolangcrawl/models"
	"time"
)

type Crawler interface {
	Crawl() ([]*models.PostRequestDto, error)
}

type CrawlService struct {
	crawlers []Crawler
}

func NewCrawlService(crawlers ...Crawler) *CrawlService {
	return &CrawlService{
		crawlers: crawlers,
	}
}

func (s *CrawlService) Crawl() ([]*models.PostRequestDto, error) {
	crawlingStart := time.Now()
	totalDtos := make([]*models.PostRequestDto, 0)
	for _, crawler := range s.crawlers {
		dtos, err := crawler.Crawl()
		if err != nil {
			fmt.Println(err)
		}
		totalDtos = append(totalDtos, dtos...)
	}
	crawlingEnd := time.Now()
	fmt.Printf("Crawling Time: %s\n", crawlingEnd.Sub(crawlingStart))

	return totalDtos, nil
}
