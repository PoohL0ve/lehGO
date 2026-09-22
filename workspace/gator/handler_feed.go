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

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.Args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUser)
	if err != nil {
		return fmt.Errorf("failed to get logged-in user: %w", err)
	}

	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       url,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to create feed: %w", err)
	}

	fmt.Printf("Feed created successfully!\n")
	fmt.Printf("ID:         %s\n", feed.ID)
	fmt.Printf("Created At: %s\n", feed.CreatedAt)
	fmt.Printf("Updated At: %s\n", feed.UpdatedAt)
	fmt.Printf("Name:       %s\n", feed.Name)
	fmt.Printf("URL:        %s\n", feed.Url)
	fmt.Printf("User ID:    %s\n", feed.UserID)

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		log.Fatalf("failed to retrieve feeds: %v", err)
	}

	if len(feeds) == 0 {
		fmt.Println("No feeds found.")
		return nil
	}

	for _, feed := range feeds {
		fmt.Printf("* Name: %s\n", feed.FeedName)
		fmt.Printf("  URL:  %s\n", feed.Url)
		fmt.Printf("  User: %s\n", feed.UserName)
		fmt.Println("---")
	}

	return nil
}
