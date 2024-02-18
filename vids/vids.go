package vids

import (
	"errors"
	"fmt"
	"log"
	"time"

	"google.golang.org/api/youtube/v3"
)

func GetPrimagenId(client *youtube.Service) (string, error) {
	maxChannels := int64(15)
	call := client.Subscriptions.List([]string{"snippet"}).Mine(true).MaxResults(maxChannels)

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
				return item.Snippet.ResourceId.ChannelId, nil
			}
		}

		nextPageToken = resp.NextPageToken

		if nextPageToken == "" {
			break
		}
	}
	return "", errors.New("Cannot find the channel Id in your channel title provide")
}

func GetUploadPlaylistId(client *youtube.Service, channelId string) (string, error) {
	callPlaylist := client.Channels.List([]string{"contentDetails"}).Id(channelId)

	resp, err := callPlaylist.Do()
	if err != nil {
		return "", err
	}

	if len(resp.Items) == 0 {
		return "", fmt.Errorf("No playlist found in the provided channel: %v", channelId)
	}

	return resp.Items[0].ContentDetails.RelatedPlaylists.Uploads, nil
}

func GetRecentVideos(client *youtube.Service, playlistId string) (map[string]string, error) {
	callVideos := client.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistId).
		MaxResults(15)

	resp, err := callVideos.Do()
	if err != nil {
		return nil, err
	}

	videoList := map[string]string{}
	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)

	for _, item := range resp.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error parsing time: %v\n", err)
			continue
		}

		if publishedAt.After(oneDayAgo) {
			fmt.Printf("%s - %s\n", item.Snippet.Title, item.Snippet.PublishedAt)
			videoList[item.Snippet.ResourceId.VideoId] = item.Snippet.Title
		}

	}
	return videoList, nil
}
