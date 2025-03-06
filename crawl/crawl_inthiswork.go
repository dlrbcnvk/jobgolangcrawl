package crawl

import (
	"fmt"
	"github.com/gocolly/colly/v2"
	"jobgolangcrawl/models"
	"strconv"
	"strings"
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
	// 지금은 bool 값으로 취소 여부를 판단하지만,
	// 채널이나 컨텍스트로 취소 여부를 판단하도록 바꿀 수 있을까? 바꿔야 할까?
	cancel := false
	c.Collector.OnHTML("div.fusion-text.fusion-text-1 > p > strong > span", func(e *colly.HTMLElement) {
		text := e.Text
		if strings.Contains(text, "죄송합니다.") {
			cancel = true
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
		dtos = append(dtos, requestDto)
	})

	c.Collector.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	page := 1
	for {
		pageStr := strconv.Itoa(page)
		pageUrl := fmt.Sprintf(c.Url, pageStr)
		err := c.Collector.Visit(pageUrl)
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
