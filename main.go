package main

import (
	"fmt"
	config "jobgolangcrawl/config"
	"jobgolangcrawl/crawl"
	"jobgolangcrawl/database"
	"jobgolangcrawl/repositories"
	"jobgolangcrawl/services"
	"log"
	"os"
	"time"
)

func main() {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "local"
	}
	config, err := config.LoadConfig(env)
	if err != nil {
		log.Fatal(err)
	}
	totalStart := time.Now()

	// 데이터베이스 초기화
	db := database.Initialize(config)

	// crawling
	crawlService := services.NewCrawlService(
		crawl.NewSaraminCrawler(),
		crawl.NewJobKoreaCrawler(),
		crawl.NewIncruitCrawler(),
		crawl.NewInThisWorkCrawler(),
	)
	dtos, err := crawlService.Crawl()
	if err != nil {
		fmt.Println(err)
		return
	}

	// create new posts
	postRepository := repositories.NewPostRepository(db)
	siteRepository := repositories.NewSiteRepository(db)
	postService := services.NewPostService(postRepository, siteRepository)
	newDtos, err := postService.CreateNewPosts(dtos)

	// send email
	mailService := services.NewMailService(newDtos, config)
	err = mailService.SendMail()
	if err != nil {
		fmt.Println(err)
		return
	}

	totalEnd := time.Now()
	fmt.Printf("Total time: %s\n", totalEnd.Sub(totalStart))
}
