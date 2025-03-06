package crawl

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"jobgolangcrawl/models"
	"strconv"
	"strings"
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
	// 인크루트는 페이지 방식이 아니라 startno 를 쿼리파라미터로 지정한다. startno부터 최대 30개를 불러온다.
	startno := 1
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
		startno++

		requestDto := models.NewPostRequestDto(postId, link, title, companyName, c.Site)
		fmt.Printf("PostRequestDto: %q\n", *requestDto)
		dtos = append(dtos, requestDto)
	})

	c.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	for {
		preStartNo := startno
		pageStr := strconv.Itoa(startno)
		err := c.Collector.Visit(c.Url + pageStr)
		if err != nil {
			fmt.Println(err)
			return dtos, err
		}

		if startno == preStartNo {
			fmt.Printf("%d 번을 끝으로 작업을 종료합니다.\n", (startno - 1))
			break
		}
	}

	fmt.Printf("=========== Crawling ends %s\n", c.Domains[0])
	return dtos, nil
}
