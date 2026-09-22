package main

import (
	"context"
	"errors"
	"fmt"
	"log"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: login <username>")
	}

	username := cmd.Args[0]

	// Look up user in the database
	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("user does not exist: %v", err)
	}

	// Update config with active user
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("User switched to: %s\n", user.Name)
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.ResetUsers(context.Background())
	if err != nil {
		log.Fatalf("failed to reset database: %v", err)
	}

	fmt.Println("Database successfully reset!")
	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("failed to retrieve users: %v", err)
	}

	for _, user := range users {
		if user.Name == s.cfg.CurrentUser {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}

	return nil
}
