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

func handlerFeedsFollow(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		return errors.New("provide url to follow")
	}
	feedUrl := cmd.Args[0]

	currentUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting currentUserDetail: %w", err)
	}

	fetchedFeed, err := s.db.GetFeedByUrl(ctx, feedUrl)
	if err == sql.ErrNoRows {
		return errors.New("feed does'nt exist")
	}

	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	now := time.Now()

	insertedFeedFollow, err := s.db.CreateFeedFollow(ctx, database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    currentUser.ID,
		FeedID:    fetchedFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error in following the feed: %w", err)
	}

	fmt.Printf("you are following the feed %s created by %s", insertedFeedFollow.FeedName, insertedFeedFollow.UserName)

	return nil
}

func handlerFeedsFollowing(s *state, cmd command) error {

	ctx := context.Background()
	currentUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting currentUserDetail: %w", err)
	}
	fetchedFollowings, err := s.db.GetAllFollowings(ctx, currentUser.ID)
	if err == sql.ErrNoRows {
		return errors.New("you are not following the feeds")
	}
	if err != nil {
		return fmt.Errorf("error getting the following: %w", err)
	}
	for _, feed := range fetchedFollowings {
		fmt.Println(feed.FeedName)
	}

	return nil
}
