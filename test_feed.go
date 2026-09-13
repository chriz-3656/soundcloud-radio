package main

import (
	"context"
	"fmt"
	"soundcloud-radio/internal/jiosaavn"
)

func main() {
	client := jiosaavn.NewClient()
	tracks, err := client.GetHomeFeed(context.Background())
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fetched %d tracks\n", len(tracks))
	if len(tracks) > 0 {
		fmt.Printf("First track: %+v\n", tracks[0])
	}
}
