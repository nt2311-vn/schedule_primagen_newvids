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

	maxChannels := int64(15)

	call := clientService.Subscriptions.List([]string{"snippet"}).Mine(true).MaxResults(maxChannels)
	nextPageToken := ""

	primagenTitles := map[string]bool{"ThePrimeTime": true}

	for {
		if nextPageToken != "" {
			call.PageToken(nextPageToken)
		}

		resp, err := call.Do()
		if err != nil {
			log.Fatalf("Error fetching Subscriptions: %v", err)
		}

		for _, item := range resp.Items {
			if _, channelExist := primagenTitles[item.Snippet.Title]; channelExist {
				fmt.Printf("Found the Primagen channel id: %s\n", item.Snippet.ResourceId.ChannelId)
				return
			}
		}

		nextPageToken = resp.NextPageToken

		if nextPageToken == "" {
			break
		}

	}
}
