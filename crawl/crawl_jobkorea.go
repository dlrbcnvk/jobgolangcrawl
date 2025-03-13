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

type JobKoreaCrawler struct {
	BaseCrawler
}

func NewJobKoreaCrawler() *JobKoreaCrawler {
	domains := make([]string, 2)
	domains[0] = "jobkorea.co.kr"
	domains[1] = "www.jobkorea.co.kr"
	return &JobKoreaCrawler{
		BaseCrawler: BaseCrawler{
			Site:      "jobkorea",
			Domains:   domains,
			Url:       "https://www.jobkorea.co.kr/Search/?stext=golang&ord=RelevanceDesc&tabType=recruit&Page_No=",
			Collector: colly.NewCollector(colly.AllowedDomains(domains...)),
		},
	}
}

func (c *JobKoreaCrawler) Crawl() ([]*models.PostRequestDto, error) {
	dtos := make([]*models.PostRequestDto, 0)
	fmt.Printf("=========== Crawling starts %s\n", c.Domains[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	c.Collector.OnHTML("section.content-recruit > article.list-empty", func(e *colly.HTMLElement) {
		cancel()
	})

	c.Collector.OnHTML("section.content-recruit.on > article.list > article.list-item", func(e *colly.HTMLElement) {
		var postId, link, title, companyName string

		// post ID
		postId = e.Attr("data-gno")
		fmt.Printf("PostId: %q\n", postId)

		// company name
		e.ForEach("div.list-section-corp > a[href]", func(_ int, e *colly.HTMLElement) {
			companyName = strings.ReplaceAll(e.Text, "\n", "")
			companyName = strings.TrimSpace(companyName)
			fmt.Printf("Company Name: %q\n", companyName)
		})

		// title
		e.ForEach("div.information-title > a[href]", func(_ int, e *colly.HTMLElement) {
			title = e.Text
			title = strings.ReplaceAll(title, "\n", "")
			title = strings.TrimSpace(title)
			fmt.Printf("Title: %q\n", title)

			// link
			link = e.Attr("href")
			link = strings.TrimSpace(link)
			if strings.Index(link, "/") == 0 {
				link = "https://www.jobkorea.co.kr" + link
			}
			fmt.Printf("Link: %q\n", link)
		})

		requestDto := models.NewPostRequestDto(postId, link, title, companyName, c.Site)
		mu.Lock()
		dtos = append(dtos, requestDto)
		mu.Unlock()
	})

	c.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	page := 1
	for {
		select {
		case <-ctx.Done():
			fmt.Println("JobKorea Crawler Done")
			break
		default:
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				pageStr := strconv.Itoa(page)
				err := c.Collector.Visit(c.Url + pageStr)
				if err != nil {
					fmt.Println(err)
					cancel()
				}
			}(page)
			page++
			time.Sleep(100 * time.Millisecond) // 작은 지연을 추가하여 요청 간의 충돌 방지
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
