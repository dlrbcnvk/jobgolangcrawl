package crawl

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"jobgolangcrawl/models"
	"strconv"
	"strings"
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
	// 지금은 bool 값으로 취소 여부를 판단하지만,
	// 채널이나 컨텍스트로 취소 여부를 판단하도록 바꿀 수 있을까? 바꿔야 할까?
	cancel := false
	c.Collector.OnHTML("section.content-recruit > article.list-empty", func(e *colly.HTMLElement) {
		cancel = true
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
		//fmt.Printf("PostRequestDto: %q\n", *requestDto)
		dtos = append(dtos, requestDto)
	})

	c.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	page := 1
	for {
		pageStr := strconv.Itoa(page)
		err := c.Collector.Visit(c.Url + pageStr)
		if err != nil {
			fmt.Println(err)
			return dtos, err
		}
		if cancel {
			fmt.Printf("%s 페이지를 끝으로 작업을 종료합니다.\n", pageStr)
			break
		}
		page++
	}

	fmt.Printf("=========== Crawling ends %s\n", c.Domains[0])
	return dtos, nil
}
