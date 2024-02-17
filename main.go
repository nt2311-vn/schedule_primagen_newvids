package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")
	fmt.Println(youtubeAPIKey)
}
