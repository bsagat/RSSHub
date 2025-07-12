package models

import (
	"fmt"
	"strings"
	"time"
)

type RSSFeed struct {
	ID        string
	CreatedAt time.Time
	Channel   Channel `xml:"channel"`
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Item        []RSSItem `xml:"item"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// To pretty print of RSSFeed struct
func (f RSSFeed) String() string {
	var sb strings.Builder

	sb.WriteString("RSS Feed:\n")
	sb.WriteString(fmt.Sprintf("  ID: %s\n", f.ID))
	sb.WriteString(fmt.Sprintf("  CreatedAt: %s\n", f.CreatedAt.Format(time.RFC1123)))
	sb.WriteString("  Channel:\n")
	sb.WriteString(fmt.Sprintf("    Title: %s\n", f.Channel.Title))
	sb.WriteString(fmt.Sprintf("    Link: %s\n", f.Channel.Link))
	sb.WriteString(fmt.Sprintf("    Description: %s\n", f.Channel.Description))
	sb.WriteString(fmt.Sprintf("    Items (%d):\n", len(f.Channel.Item)))

	for i, item := range f.Channel.Item {
		sb.WriteString(fmt.Sprintf("      Item #%d:\n", i+1))
		sb.WriteString(fmt.Sprintf("        Title: %s\n", item.Title))
		sb.WriteString(fmt.Sprintf("        Link: %s\n", item.Link))
		sb.WriteString(fmt.Sprintf("        Description: %s\n", item.Description))
		sb.WriteString(fmt.Sprintf("        PubDate: %s\n", item.PubDate))
	}

	return sb.String()
}
