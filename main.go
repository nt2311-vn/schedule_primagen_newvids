package main

import (
	"fmt"
	"log"

	"github.com/nt2311-vn/schedule_primagen_newvids/auth"
	"github.com/nt2311-vn/schedule_primagen_newvids/vids"
)

func main() {
	mapVids, err := getVideos()
	if err != nil {
		log.Fatalf("Error on get videos from the channel: %v", err)
	}

	fmt.Printf("Found %d recent video(s)\n", len(mapVids))

	if len(mapVids) > 0 {
		for key, value := range mapVids {
			fmt.Printf("Title: %v, video info: %v\n", *value, key)
		}
	}
}

func getVideos() (map[string]*vids.VideoInfo, error) {
	establishClient, errClient := auth.GetYoutubeService()

	if errClient != nil {
		return nil, errClient
	}

	channelId, errGetId := vids.GetPrimagenId(establishClient)

	if errGetId != nil {
		return nil, errGetId
	}

	playlistId, errGetPlaylist := vids.GetUploadPlaylistId(establishClient, channelId)

	if errGetPlaylist != nil {
		return nil, errGetPlaylist
	}

	recentPlaylists, errGetVideos := vids.GetRecentVideos(establishClient, playlistId)

	if errGetVideos != nil {
		return nil, errGetVideos
	}
	return recentPlaylists, nil
}
