package main

import (
	"fmt"
	"log"

	"github.com/nt2311-vn/schedule_primagen_newvids/auth"
)

func main() {
	clientService, errGetService := auth.GetAuth()

	if errGetService != nil {
		log.Fatalf("Error at getting service from api youtube: %v", errGetService)
	}

	call := clientService.Subscriptions.List([]string{"snippet"}).Mine(true)

	resp, err := call.Do()
	if err != nil {
		log.Fatalf("Error fetching Subscriptions: %v", err)
	}

	for _, item := range resp.Items {
		fmt.Printf("%s : %s\n", item.Snippet.ResourceId.ChannelId, item.Snippet.Title)
	}
}
