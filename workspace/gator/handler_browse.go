package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/database"
)

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := int32(2)
	if len(cmd.Args) > 0 {
		parsedLimit, err := strconv.Atoi(cmd.Args[0])
		if err != nil || parsedLimit <= 0 {
			return errors.New("limit must be a positive integer")
		}
		limit = int32(parsedLimit)
	}

	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  limit,
	})
	if err != nil {
		return fmt.Errorf("failed to get posts: %w", err)
	}

	if len(posts) == 0 {
		fmt.Println("No posts found. Make sure you are following feeds and running the aggregator.")
		return nil
	}

	fmt.Printf("--- Browsing Latest %d Posts for %s ---\n\n", len(posts), user.Name)
	for _, post := range posts {
		pubStr := "Unknown date"
		if post.PublishedAt.Valid {
			pubStr = post.PublishedAt.Time.Format("Jan 02, 2006")
		}

		fmt.Printf("Title: %s\n", post.Title)
		fmt.Printf("Feed:  %s\n", post.FeedName)
		fmt.Printf("Date:  %s\n", pubStr)
		fmt.Printf("URL:   %s\n", post.Url)
		if post.Description.Valid && post.Description.String != "" {
			fmt.Printf("Summary: %s\n", post.Description.String)
		}
		fmt.Println(strings.Repeat("-", 40))
	}

	return nil
}
