package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"time"
)

func main() {
	url := "https://i1.sndcdn.com/artworks-000247656247-4933cw-large.jpg"

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println("GET ERROR:", err)
		return
	}
	defer resp.Body.Close()

    fmt.Println("Status:", resp.StatusCode)

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		fmt.Println("DECODE ERROR:", err)
		return
	}
    fmt.Println("Success! Bounds:", img.Bounds())
}
