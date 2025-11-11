package main

import (
	"errors"
	"fmt"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return errors.New("username not provided")
	}

	userName := cmd.Args[0]

	err := s.cfg.SetUser(userName)
	if err != nil {
		return err
	}

	fmt.Printf("the user %s is successfully set.\n", userName)

	return nil
}
