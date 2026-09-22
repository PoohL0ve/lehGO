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

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 2 {
		return errors.New("usage: addfeed <name> <url>")
	}

	name := cmd.Args[0]
	url := cmd.Args[1]

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

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("feed created, but failed to follow: %w", err)
	}

	fmt.Printf("Feed created successfully!\n")
	fmt.Printf("ID:         %s\n", feed.ID)
	fmt.Printf("Name:       %s\n", feed.Name)
	fmt.Printf("URL:        %s\n", feed.Url)
	fmt.Printf("User:       %s\n", user.Name)
	fmt.Printf("Now following feed: %s\n", ffRow.FeedName)

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

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: follow <url>")
	}

	url := cmd.Args[0]

	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found for URL %s: %w", url, err)
	}

	ffRow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feed.ID,
	})
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("User %s is now following feed: %s\n", ffRow.UserName, ffRow.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, user database.User) error {
	feedFollows, err := s.db.GetFeedFollowsForUser(context.Background(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to get feed follows: %w", err)
	}

	if len(feedFollows) == 0 {
		fmt.Println("You are not following any feeds.")
		return nil
	}

	fmt.Printf("Feeds followed by %s:\n", user.Name)
	for _, ff := range feedFollows {
		fmt.Printf("* %s\n", ff.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: unfollow <url>")
	}

	url := cmd.Args[0]

	// 1. Get the feed by URL
	feed, err := s.db.GetFeedByURL(context.Background(), url)
	if err != nil {
		return fmt.Errorf("feed not found for URL %s: %w", url, err)
	}

	// 2. Delete the feed follow record matching user ID and feed ID
	err = s.db.DeleteFeedFollow(context.Background(), database.DeleteFeedFollowParams{
		UserID: user.ID,
		FeedID: feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to unfollow feed: %w", err)
	}

	fmt.Printf("%s has unfollowed feed: %s\n", user.Name, feed.Name)
	return nil
}
