package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/database"
	"github.com/google/uuid"
)

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: register <name>")
	}

	name := cmd.Args[0]

	newUserParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
	}

	user, err := s.db.CreateUser(context.Background(), newUserParams)
	if err != nil {
		log.Fatalf("could not create user: %v", err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("could not set current user: %w", err)
	}

	fmt.Printf("User created successfully: %s\n", user.Name)
	fmt.Printf("User Debug Data: %+v\n", user)

	return nil
}
