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

// PrintArticleList prints a formatted list of articles for a specific feed.
func PrintArticleList(articles []*models.RSSItem, feedName string) {
	format := `%d. [%s] %s
   %s

`

	fmt.Printf("# Feed: %s\n\n", feedName)
	for i, article := range articles {
		fmt.Printf(format, i+1, article.PubDate, article.Title, article.Link)
	}
}

// PrintFeedsList prints a formatted list of available RSS feeds to the console.
func PrintFeedsList(feeds []*models.Feed) {
	format := `%d. Name: %s
   URL: %s
   Added: %s

`

	fmt.Print("# Available RSS Feeds\n\n")
	for i, feed := range feeds {
		fmt.Printf(format, i+1, feed.Name, feed.URL, feed.CreatedAt.Format(time.DateTime))
	}
}
