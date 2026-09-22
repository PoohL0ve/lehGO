package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/database"
	"github.com/PoohL0ve/lehGO/workspace/gator/internal/rss"
	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) < 1 {
		return errors.New("usage: agg <time_between_reqs>")
	}

	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("invalid duration format: %w", err)
	}

	fmt.Printf("Collecting feeds every %s...\n", timeBetweenReqs)

	ticker := time.NewTicker(timeBetweenReqs)
	defer ticker.Stop()

	// Executes immediately on start, then wait on ticker ticks
	for ; ; <-ticker.C {
		if err := scrapeFeeds(s); err != nil {
			log.Printf("Error scraping feed: %v", err)
		}
	}
}

func parsePublishedAt(pubDate string) sql.NullTime {
	if pubDate == "" {
		return sql.NullTime{Valid: false}
	}

	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		time.RFC3339,
		"Mon, 02 Jan 2006 15:04:05 -0700",
		"Mon, 02 Jan 2006 15:04:05 MST",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, pubDate); err == nil {
			return sql.NullTime{Time: t.UTC(), Valid: true}
		}
	}

	return sql.NullTime{Valid: false}
}

func scrapeFeeds(s *state) error {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get next feed: %w", err)
	}

	now := time.Now().UTC()
	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFetchedAt: sql.NullTime{Time: now, Valid: true},
		ID:            feed.ID,
	})
	if err != nil {
		return fmt.Errorf("failed to mark feed fetched: %w", err)
	}

	rssFeed, err := rss.FetchFeed(context.Background(), feed.Url)
	if err != nil {
		return fmt.Errorf("failed to fetch RSS feed %s: %w", feed.Url, err)
	}

	for _, item := range rssFeed.Channel.Item {
		description := sql.NullString{
			String: item.Description,
			Valid:  item.Description != "",
		}

		pubTime := parsePublishedAt(item.PubDate)

		_, err := s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   now,
			UpdatedAt:   now,
			Title:       item.Title,
			Url:         item.Link,
			Description: description,
			PublishedAt: pubTime,
			FeedID:      feed.ID,
		})

		if err != nil {
			// Ignore duplicate URL unique constraint errors
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") ||
				strings.Contains(err.Error(), "posts_url_key") {
				continue
			}
			log.Printf("Failed to save post '%s': %v", item.Title, err)
		}
	}

	fmt.Printf("Fetched %s: Saved posts to database\n", feed.Name)
	return nil
}
