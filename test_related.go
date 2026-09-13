package main

import (
	"context"
	"fmt"
	"soundcloud-radio/internal/jiosaavn"
)

func main() {
	client := jiosaavn.NewClient()
	tracks, err := client.GetRelatedTracks(context.Background(), "jx2G3Mjl", 10)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("Fetched %d related tracks\n", len(tracks))
	if len(tracks) > 0 {
		fmt.Printf("First related track: %+v\n", tracks[0])
	}
}
