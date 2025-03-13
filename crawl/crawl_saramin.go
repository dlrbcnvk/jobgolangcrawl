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

type SaraminCrawler struct {
	BaseCrawler
}

func NewSaraminCrawler() *SaraminCrawler {
	domains := make([]string, 2)
	domains[0] = "saramin.co.kr"
	domains[1] = "www.saramin.co.kr"
	return &SaraminCrawler{
		BaseCrawler: BaseCrawler{
			Site:      "saramin",
			Domains:   domains,
			Url:       "https://www.saramin.co.kr/zf_user/jobs/list/job-category?cat_kewd=223&sort=RD&page=",
			Collector: colly.NewCollector(colly.AllowedDomains(domains...)),
		},
	}
}

func (c *SaraminCrawler) Crawl() ([]*models.PostRequestDto, error) {
	dtos := make([]*models.PostRequestDto, 0)
	fmt.Printf("=========== Crawling starts %s\n", c.Domains[0])

	var wg sync.WaitGroup
	var mu sync.Mutex
	ctx, cancel := context.WithCancel(context.Background())
	c.Collector.OnHTML("div.info_empty", func(e *colly.HTMLElement) {
		cancel()
	})

	c.Collector.OnHTML("div.box_item", func(e *colly.HTMLElement) {
		var postId, link, title, companyName string

		// company name
		e.ForEach("div.col.company_nm > a[href]", func(_ int, e *colly.HTMLElement) {
			companyName = strings.ReplaceAll(e.Text, "\n", "")
			companyName = strings.TrimSpace(companyName)
			fmt.Printf("Company Name: %q\n", companyName)
		})

		e.ForEach("div.col.notification_info > div.job_tit > a.str_tit[href]", func(_ int, e *colly.HTMLElement) {
			// link
			hrefAttr := e.Attr("href")
			link = "https://" + c.Domains[0] + hrefAttr
			fmt.Printf("Link found: %q -> %s\n", e.Text, link)

			// post_id
			idAttr := e.Attr("id")
			postId = strings.ReplaceAll(idAttr, "rec_link_", "")
			fmt.Printf("Post Id: %q\n", postId)

			// title
			title = e.Text
			fmt.Printf("Title: %q\n", title)

		})

		// PostRequestDto
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
			fmt.Println("Saramin Crawler done")
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
