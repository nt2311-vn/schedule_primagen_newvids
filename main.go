package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func main() {
	errLoadEnv := godotenv.Load()
	if errLoadEnv != nil {
		log.Fatal(errLoadEnv)
	}

	fmt.Println("Load environment varialbe comple")
	youtubeAPIKey := os.Getenv("YOUTUBE_API_KEY")

	ctx := context.Background()

	client, errConnection := youtube.NewService(ctx, option.WithAPIKey(youtubeAPIKey))

	if errConnection != nil {
		log.Fatalf("Error establishing youtube client %v", errConnection)
	}

	call := client.Subscriptions.List([]string{"snippet"}).Mine(true)

	subscribeList, errGetList := call.Do()

	if errGetList != nil {
		log.Fatalf("Error get subscription list %v", errGetList)
	}

	fmt.Println("List of subscription channel")

	for _, item := range subscribeList.Items {
		fmt.Printf("ID: %s, Title: %s\n", item.Snippet.ResourceId.ChannelId, item.Snippet.Title)
	}
}
