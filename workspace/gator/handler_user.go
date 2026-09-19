package main

import "fmt"

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	username := cmd.Args[0]
	err := s.cfg.SetUser(username)
	if err != nil {
		return fmt.Errorf("could not set user: %w", err)
	}

	fmt.Printf("%s user has been set\n", username)
	return nil
}
