package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/almuiz021/gatorcli/internal/database"
)

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(s *state, cmd command) error {

	return func(s *state, cmd command) error {
		ctx := context.Background()
		currentUser, err := s.db.GetUser(ctx, s.cfg.CurrentUserName)
		if err == sql.ErrNoRows {
			return errors.New("no user found ")
		}
		if err != nil {
			return fmt.Errorf("error getting currentUserDetail: %w", err)
		}

		return handler(s, cmd, currentUser)
	}
}
