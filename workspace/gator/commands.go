package main

import "errors"

type command struct {
	Name string
	Args []string
}

// holds all commands for the CLI and their handler functions
type commands struct {
	registeredCommands map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.registeredCommands[cmd.Name]
	if !ok {
		errors.New("Command not found")
	}

	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.registeredCommands[name] = f
}
