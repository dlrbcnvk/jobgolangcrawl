package crawl

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"jobgolangcrawl/models"
	"strconv"
	"strings"
	"sync"
	"time"
)

type IncruitCrawler struct {
	BaseCrawler
}

func NewIncruitCrawler() *IncruitCrawler {
	domains := make([]string, 2)
	domains[0] = "incruit.com"
	domains[1] = "search.incruit.com"
	return &IncruitCrawler{
		BaseCrawler: BaseCrawler{
			Site:      "incruit",
			Domains:   domains,
			Url:       "https://search.incruit.com/list/search.asp?col=job&kw=golang&startno=",
			Collector: colly.NewCollector(colly.AllowedDomains(domains...)),
		},
	}
}

func (c *IncruitCrawler) Crawl() ([]*models.PostRequestDto, error) {
	dtos := make([]*models.PostRequestDto, 0)
	fmt.Printf("=========== Crawling starts %s\n", c.Domains[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())

	c.Collector.OnHTML("body", func(e *colly.HTMLElement) {
		contentLength := e.DOM.Find("ul.c_row").Length()
		if contentLength == 0 {
			cancel()
		}
	})

	c.Collector.OnHTML("ul.c_row", func(e *colly.HTMLElement) {
		var postId, link, title, companyName string

		// postId
		postId = e.Attr("jobno")

		// company name
		e.ForEach("div.cell_first > div.cl_top > a[href]", func(_ int, e *colly.HTMLElement) {
			companyName = e.Text
			companyName = strings.TrimSpace(companyName)
		})

		e.ForEach("div.cell_mid > div.cl_top > a[href]", func(_ int, e *colly.HTMLElement) {
			// title
			title = e.Text

			// link
			link = e.Attr("href")
		})

		requestDto := models.NewPostRequestDto(postId, link, title, companyName, c.Site)
		mu.Lock()
		dtos = append(dtos, requestDto)
		mu.Unlock()
	})

	c.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// 인크루트는 페이지 방식이 아니라 startno 를 쿼리파라미터로 지정한다. startno부터 최대 30개를 불러온다.
	startno := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("Incruit Crawler Done")
			break
		default:
			wg.Add(1)
			go func(startno int) {
				defer wg.Done()
				pageStr := strconv.Itoa(startno)
				err := c.Collector.Visit(c.Url + pageStr)
				if err != nil {
					fmt.Println(err)
					cancel()
				}
			}(startno)
			startno += 30
			time.Sleep(100 * time.Millisecond)
		}
		if ctx.Err() != nil {
			break
		}
	}
	wg.Wait()
	cancel()

	fmt.Printf("=========== Crawling ends %s\n", c.Domains[0])
	return dtos, nil
}
