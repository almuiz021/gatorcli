package main

import (
	"errors"
)

type command struct {
	Name string
	Args []string
}

type commands struct {
	handlers map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	cmdToRun, ok := c.handlers[cmd.Name]
	if !ok {
		return errors.New("command not found")
	}

	return cmdToRun(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}
