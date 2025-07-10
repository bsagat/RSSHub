package utils

import (
	"RSSHub/internal/domain/models"
	"fmt"
	"time"
)

func PrintHelp() {
	text := `./rsshub --help

  Usage:
    rsshub COMMAND [OPTIONS]

  Common Commands:
       add             add new RSS feed
       set-interval    set RSS fetch interval
       set-workers     set number of workers
       list            list available RSS feeds
       delete          delete RSS feed
       articles        show latest articles
       fetch           starts the background process that periodically fetches and processes RSS feeds using a worker pool
`
	fmt.Println(text)
}

// Feed: tech-crunch

// 1. [2025-06-18] Apple announces new M4 chips for MacBook Pro
//    https://techcrunch.com/apple-announces-m4/

// 2. [2025-06-17] OpenAI launches GPT-5 with multimodal capabilities
//    https://techcrunch.com/openai-launches-gpt-5/

// 3. [2025-06-16] Google unveils new privacy tools at I/O 2025
//    https://techcrunch.com/google-privacy-io-2025/

// 4. [2025-06-15] TikTok introduces developer platform for integrations
//    https://techcrunch.com/tiktok-developer-platform/

// 5. [2025-06-14] Microsoft Teams gets AI-powered meeting summarization
//    https://techcrunch.com/microsoft-teams-ai-summary/

func PrintArticleList(articles []models.RSSItem, feedName string) {
	format :=
		`%d. [%s] %s
	 %s 
	 
	 `

	fmt.Printf("Feed: %s\n\n", feedName)
	for i, article := range articles {
		fmt.Printf(format, i+1, article.PubDate, article.Title, article.Link)
	}
}

func PrintFeedsList(feeds []models.RSSFeed) {
	format := `%d. Name: %s 
	 URL: %s 
	 Added: %s 
	 
	 `

	fmt.Print("# Available RSS Feeds\n\n")
	for i, feed := range feeds {
		fmt.Printf(format, i+1, feed.Channel.Title, feed.Channel.Link, feed.CreatedAt.Format(time.DateTime))
	}
}
