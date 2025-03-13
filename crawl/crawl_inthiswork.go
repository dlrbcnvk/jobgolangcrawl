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

type InThisWorkCrawler struct {
	BaseCrawler
}

func NewInThisWorkCrawler() *InThisWorkCrawler {
	domains := make([]string, 2)
	domains[0] = "inthiswork.com"
	return &InThisWorkCrawler{
		BaseCrawler: BaseCrawler{
			Site:      "inthiswork",
			Domains:   domains,
			Url:       "https://inthiswork.com/page/%s?s=golang",
			Collector: colly.NewCollector(colly.AllowedDomains(domains...)),
		},
	}
}

func (c *InThisWorkCrawler) Crawl() ([]*models.PostRequestDto, error) {
	dtos := make([]*models.PostRequestDto, 0)
	fmt.Printf("=========== Crawling starts %s\n", c.Domains[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	c.Collector.OnHTML("div.fusion-text.fusion-text-1 > p > strong > span", func(e *colly.HTMLElement) {
		text := e.Text
		if strings.Contains(text, "죄송합니다.") {
			cancel()
		}
	})

	c.Collector.OnHTML("div.fusion-posts-container.fusion-posts-container-pagination > article > div > div > div > h2.blog-shortcode-post-title.entry-title > a[href]", func(e *colly.HTMLElement) {
		var postId, link, title, companyName string

		// link
		link = e.Attr("href")
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)

		// postId
		split := strings.Split(link, "/")
		postId = split[len(split)-1]

		// companyName, title
		text := e.Text
		split = strings.Split(text, "｜")
		title = split[0]
		companyName = split[1]

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
			fmt.Println("InThisWork Crawler Done")
			break
		default:
			wg.Add(1)
			go func(page int) {
				defer wg.Done()
				pageStr := strconv.Itoa(page)
				pageUrl := fmt.Sprintf(c.Url, pageStr)
				err := c.Collector.Visit(pageUrl)
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
