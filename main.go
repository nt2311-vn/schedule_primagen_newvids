package main

import (
	"log"

	"github.com/nt2311-vn/schedule_primagen_newvids/auth"
)

func main() {
	establishClient, errClient := auth.GetYoutubeService()

	if errClient != nil {
		log.Fatalf("Cannot connect to youtube service", errClient)
	}
}
