package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/almuiz021/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerGetRequests(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("expected time")
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}

	fmt.Printf("Collecting feeds every %v\n\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	// run immediately, then on each tick
	for first := true; ; first = false {
		if !first {
			<-ticker.C
		}
		if err := scrapeFeeds(s); err != nil {
			// If GetNextFeedToFetch returns sql.ErrNoRows, just continue.
			if errors.Is(err, sql.ErrNoRows) {
				continue
			}
			log.Printf("scrape error: %v", err)
		}
	}
}
func scrapeFeeds(s *state) error {
	ctx := context.Background()

	feed, err := s.db.GetNextFeedToFetch(ctx)
	if err != nil {
		// bubble up so caller can handle sql.ErrNoRows
		return err
	}

	marked, err := s.db.MarkFeedFetched(ctx, feed.ID)
	if err != nil {
		return fmt.Errorf("mark fetch: %w", err)
	}

	rss, err := s.fetchFeed(ctx, marked.Url)
	if err != nil {
		return fmt.Errorf("fetch feed: %w", err)
	}

	log.Printf("Feed %s collected, %d posts found\n\n\n", marked.Name, len(rss.Channel.Item))

	for _, item := range rss.Channel.Item {
		log.Printf("📝  %s\n🔗  %s\n\n", item.Title, item.Link)

		now := time.Now().UTC()

		var publishedAt sql.NullTime
		if t, err := timeParser(item.PubDate); err == nil {
			publishedAt = sql.NullTime{Time: t, Valid: true}
		} else {
			log.Printf("warning: could not parse time %q: %v", item.PubDate, err)
		}

		_, err := s.db.CreatePost(ctx, database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: now,
			UpdatedAt: now,
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: publishedAt,
			FeedID:      marked.ID,
		})
		if err != nil {
			log.Printf("error creating post for %q: %v", item.Title, err)
			continue
		}
	}
	fmt.Printf("\n\n--------------------------------------\n\n\n")

	return nil
}

func timeParser(pubDate string) (time.Time, error) {
	if pubDate == "" {
		return time.Time{}, errors.New("pubDate is empty")
	}

	layouts := []string{
		time.RFC1123Z, // "Mon, 02 Jan 2006 15:04:05 -0700"
		time.RFC1123,  // "Mon, 02 Jan 2006 15:04:05 MST"
		time.RFC822Z,  // "02 Jan 06 15:04 -0700"
		time.RFC822,   // "02 Jan 06 15:04 MST"
		time.RFC3339,  // ISO8601 feeds
		time.RFC3339Nano,
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, pubDate); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unrecognized time format: %q", pubDate)
}
