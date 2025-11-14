package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Description string    `xml:"description"`
		Link        string    `xml:"link"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	PubDate     string `xml:"pubDate"`
	Description string `xml:"description"`
	// GUID        string `xml:"gUID"`
}

func (s *state) fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	rssFeed := RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &rssFeed, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	res, err := s.httpClient.Do(req)
	if err != nil {
		return &rssFeed, fmt.Errorf("error responding request: %w", err)
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return &rssFeed, fmt.Errorf("error reading body: %w", err)
	}
	defer res.Body.Close()

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return &rssFeed, fmt.Errorf("error unmarshalling data: %w", err)
	}

	rssFeed.Channel.Title = html.UnescapeString(rssFeed.Channel.Title)
	rssFeed.Channel.Description = html.UnescapeString(rssFeed.Channel.Description)
	for i := range rssFeed.Channel.Item {
		rssFeed.Channel.Item[i].Title = html.UnescapeString(rssFeed.Channel.Item[i].Title)
		rssFeed.Channel.Item[i].Description = html.UnescapeString(rssFeed.Channel.Item[i].Description)
	}

	return &rssFeed, nil

}
