package utils

import (
	"RSSHub/internal/domain/models"
	"fmt"
	"strings"
	"time"
)

func PrintHelp() {
	text := `
  Usage:
    rsshub COMMAND [OPTIONS]

  Common Commands:
       add             add new RSS feed
       status          show current status of application
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

// PrettyDuration returns string information about duration in pretty format
// Examples:
// 15s                  => "15 seconds"
// 1m20s                => "1 minute 20 seconds"
// 1h15m0s              => "1 hour 15 minutes"
// 1h30m0s              => "1 hour 30 minutes"
// 2m3s                 => "2 minutes 3 seconds"
// 1h1m12s              => "1 hour 1 minute 12 seconds"
// 0s                   => "0 seconds"
func PrettyDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	var parts []string

	if hours > 0 {
		unit := "hour"
		if hours > 1 {
			unit = "hours"
		}
		parts = append(parts, fmt.Sprintf("%d %s", hours, unit))
	}

	if minutes > 0 {
		unit := "minute"
		if minutes > 1 {
			unit = "minutes"
		}
		parts = append(parts, fmt.Sprintf("%d %s", minutes, unit))
	}

	if seconds > 0 || len(parts) == 0 {
		unit := "second"
		if seconds != 1 {
			unit = "seconds"
		}
		parts = append(parts, fmt.Sprintf("%d %s", seconds, unit))
	}

	return joinParts(parts)
}

func PrettyRssConfig(c *models.RssConfig) string {
	var sb strings.Builder

	sb.WriteString("Current configuration:\n")
	status := "stopped"
	if c.Run {
		status = "running"
	}

	sb.WriteString(fmt.Sprintf("  Status:         %v\n", status))
	sb.WriteString(fmt.Sprintf("  Worker count:   %d\n", c.WorkerCount))
	sb.WriteString(fmt.Sprintf("  Timer interval: %s", PrettyDuration(c.TimerInterval)))

	return sb.String()
}

func joinParts(parts []string) string {
	switch len(parts) {
	case 0:
		return ""
	case 1:
		return parts[0]
	case 2:
		return parts[0] + " " + parts[1]
	default:
		return parts[0] + " " + parts[1] + " " + parts[2]
	}
}
