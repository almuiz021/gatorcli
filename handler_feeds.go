package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/almuiz021/gatorcli/internal/database"
	"github.com/google/uuid"
)

func handlerAddFeed(s *state, cmd command, user database.User) error {

	ctx := context.Background()

	if len(cmd.Args) < 1 {
		return errors.New("expected both feed name and url")
	}

	if len(cmd.Args) < 2 {
		return errors.New("expected both feed name and url")
	}

	feedName := cmd.Args[0]
	feedUrl := cmd.Args[1]

	now := time.Now()
	insertedFeed, err := s.db.CreateFeed(ctx, database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      feedName,
		Url:       feedUrl,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error inserting feed: %w", err)
	}

	addFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    user.ID,
		FeedID:    insertedFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error following feed: %w", err)
	}
	fmt.Printf("you are now following : %s", addFollow.FeedName)

	return nil
}

func handlerFeeds(s *state, cmd command) error {

	ctx := context.Background()

	fetchedAllFeeds, err := s.db.GetFeedByUserName(ctx)
	if err == sql.ErrNoRows {
		return fmt.Errorf("no feeds to fetch: %w", err)
	}
	if err != nil {
		return fmt.Errorf("error in getting feeds: %w", err)
	}

	for _, feed := range fetchedAllFeeds {
		fmt.Printf("- %s\n- %s\n", feed.Name, feed.FeederName)
	}

	return nil
}
