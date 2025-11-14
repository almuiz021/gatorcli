package main

import (
	"context"
	"fmt"
)

func handlerGetRequests(s *state, cmd command) error {

	ctx := context.Background()
	rssFeed, err := s.fetchFeed(ctx, s.baseUrl)
	if err != nil {
		return err
	}
	fmt.Printf("Feed: %+v\n", rssFeed)

	return nil
}
