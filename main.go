package main

import (
	"fmt"
	"log"

	"github.com/nt2311-vn/schedule_primagen_newvids/auth"
	"github.com/nt2311-vn/schedule_primagen_newvids/vids"
)

func main() {
	establishClient, errClient := auth.GetYoutubeService()

	if errClient != nil {
		log.Fatalf("Cannot connect to youtube service: %v", errClient)
	}

	channelId, errGetId := vids.GetPrimagenId(establishClient)

	if errGetId != nil {
		log.Fatalf("Cannot get id from the primagen channel: %v", errGetId)
	}

	playlistId, errGetPlaylist := vids.GetUploadPlaylistId(establishClient, channelId)

	if errGetPlaylist != nil {
		log.Fatalf("Cannot trace playlist id on channel: %v, error:%v", channelId, errGetPlaylist)
	}

	recentPlaylists, errGetVideos := vids.GetRecentVideos(establishClient, playlistId)

	if errGetVideos != nil {
		log.Fatalf("Error on get recent videos: %v", recentPlaylists)
	}

	fmt.Printf("The recent video list: %v", recentPlaylists)
}
