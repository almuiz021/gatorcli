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

func handlerLogin(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) == 0 {
		return errors.New("username not provided")
	}

	userName := cmd.Args[0]

	u, err := s.db.GetUser(ctx, userName)
	if err == sql.ErrNoRows {
		return fmt.Errorf("user %s doesn't exists", userName)
	}
	if err != nil {
		return fmt.Errorf("error getting user: %w", err)
	}

	if err := s.cfg.SetUser(u.Name); err != nil {
		return fmt.Errorf("set user: %w", err)
	}
	fmt.Printf("the user %s is successfully set.\n", userName)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	ctx := context.Background()
	if len(cmd.Args) == 0 {
		return errors.New("username to register required")
	}

	now := time.Now()
	userName := cmd.Args[0]

	fetchedUser, err := s.db.GetUser(ctx, userName)
	if err == sql.ErrNoRows {
		insertedUser, err := s.db.CreateUser(ctx, database.CreateUserParams{
			ID:        uuid.New(),
			Name:      userName,
			CreatedAt: now,
			UpdatedAt: now,
		})
		if err != nil {
			return fmt.Errorf("error in insertingUser: %v", err)
		}

		if err := s.cfg.SetUser(userName); err != nil {
			return fmt.Errorf("set user: %w", err)
		}

		fmt.Printf("user %s created\n", userName)
		fmt.Println(insertedUser)
		return nil
	}
	if err != nil {
		return fmt.Errorf("error in fetchingUser: %v", err)
	}

	return fmt.Errorf("user %s already exists", fetchedUser.Name)
}

func handlerReset(s *state, cmd command) error {

	ctx := context.Background()
	if err := s.db.DeleteAllUsers(ctx); err != nil {
		return fmt.Errorf("error deleting all users: %w", err)
	}
	return nil
}
