package vids

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"time"

	"google.golang.org/api/youtube/v3"
)

type Video struct {
	Title     string
	LengthMin int
}

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

func iso8601ToMinutes(duration string) (int, error) {
	totalMininutes := 0

	re := regexp.MustCompile(`PT(\d+H)?(\d+M)?(\d+S)?`)

	matches := re.FindStringSubmatch(duration)

	if matches == nil {
		return 0, fmt.Errorf("Invalid 8601 duration format: %v", duration)
	}

	if matches[1] != "" {
		hours, err := strconv.Atoi(matches[1][:len(matches[1])-1])
		if err != nil {
			return 0, err
		}

		totalMininutes += hours * 60
	}

	if matches[2] != "" {
		minutes, err := strconv.Atoi(matches[2][:len(matches[2])-1])
		if err != nil {
			return 0, err
		}

		totalMininutes += minutes
	}

	if matches[3] != "" {
		seconds, err := strconv.Atoi(matches[3][:len(matches[3])-1])
		if err != nil {
			return 0, err
		}

		totalMininutes += seconds / 60
	}

	return totalMininutes, nil
}

func GetRecentVideos(client *youtube.Service, playlistId string) (map[string]Video, error) {
	callVideos := client.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistId).
		MaxResults(15)

	resp, err := callVideos.Do()
	if err != nil {
		return nil, err
	}

	videoList := map[string]Video{}
	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)

	for _, item := range resp.Items {
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error parsing time: %v\n", err)
			continue
		}

		if publishedAt.After(oneDayAgo) {

			toMinutes, errCovertMins := iso8601ToMinutes(item.ContentDetails.EndAt)

			if errCovertMins != nil {
				log.Fatalf("Cannot covert video length string to minutes: %v", errCovertMins)
			}

			videoList[item.Snippet.ResourceId.VideoId] = Video{item.Snippet.Title, toMinutes}
		}

	}
	return videoList, nil
}
