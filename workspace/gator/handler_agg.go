package main

import (
	"context"
	"fmt"
	"log"

	"github.com/PoohL0ve/lehGO/workspace/gator/internal/rss"
)

func handlerAgg(s *state, cmd command) error {
	feed, err := rss.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		log.Fatalf("failed to fetch feed: %v", err)
	}

	fmt.Printf("%+v\n", feed)
	return nil
}
