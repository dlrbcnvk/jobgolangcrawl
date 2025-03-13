package services

import (
	"fmt"
	"jobgolangcrawl/models"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(len(s.crawlers))
	var mutex = &sync.Mutex{}
	for _, crawler := range s.crawlers {
		go func() {
			defer wg.Done()
			dtos, err := crawler.Crawl()
			if err != nil {
				fmt.Println(err)
			}
			mutex.Lock()
			totalDtos = append(totalDtos, dtos...)
			mutex.Unlock()
		}()
	}
	wg.Wait()
	crawlingEnd := time.Now()
	fmt.Printf("Crawling Time: %s\n", crawlingEnd.Sub(crawlingStart))

	return totalDtos, nil
}
