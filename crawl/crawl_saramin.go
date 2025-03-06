package crawl

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"jobgolangcrawl/models"
	"strconv"
	"strings"
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
	// 지금은 bool 값으로 취소 여부를 판단하지만,
	// 채널이나 컨텍스트로 취소 여부를 판단하도록 바꿀 수 있을까? 바꿔야 할까?
	cancel := false
	c.Collector.OnHTML("div.info_empty", func(e *colly.HTMLElement) {
		cancel = true
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
