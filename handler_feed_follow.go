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

func handlerFeedsFollow(s *state, cmd command, user database.User) error {
	ctx := context.Background()
	if len(cmd.Args) < 1 {
		return errors.New("provide url to follow")
	}
	feedUrl := cmd.Args[0]

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
		UserID:    user.ID,
		FeedID:    fetchedFeed.ID,
	})
	if err != nil {
		return fmt.Errorf("error in following the feed: %w", err)
	}

	fmt.Printf("you are following the feed %s created by %s", insertedFeedFollow.FeedName, insertedFeedFollow.UserName)

	return nil
}

func handlerFeedsFollowing(s *state, cmd command, user database.User) error {

	ctx := context.Background()
	fetchedFollowings, err := s.db.GetAllFollowings(ctx, user.ID)
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

func handlerFeedsUnFollow(s *state, cmd command, user database.User) error {

	feedUrl := cmd.Args[0]

	ctx := context.Background()
	feedDetails, err := s.db.GetFeedByUrl(ctx, feedUrl)
	if err == sql.ErrNoRows {
		return errors.New("feed does'nt exist")
	}
	if err != nil {
		return fmt.Errorf("error getting feed by url: %w", err)
	}

	if err := s.db.DeleteFeedFollow(ctx, database.DeleteFeedFollowParams{UserID: user.ID, FeedID: feedDetails.ID}); err != nil {
		return errors.New("error deleting the feed")
	}

	return nil
}
