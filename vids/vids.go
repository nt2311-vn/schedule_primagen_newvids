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

type VideoInfo struct {
	Title      string
	LengthMins int
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

func GetRecentVideos(client *youtube.Service, playlistId string) (map[string]*VideoInfo, error) {
	callPlaylists := client.PlaylistItems.List([]string{"snippet"}).
		PlaylistId(playlistId).
		MaxResults(15)

	resp, err := callPlaylists.Do()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	oneDayAgo := now.AddDate(0, 0, -1)
	videoList := make(map[string]*VideoInfo)

	for _, item := range resp.Items {
		if item.Snippet == nil || item.Snippet.ResourceId == nil {
			log.Println("Skipping item due to nil snipper or resource id")
			continue
		}
		publishedAt, err := time.Parse(time.RFC3339, item.Snippet.PublishedAt)
		if err != nil {
			log.Printf("Error parsing time: %v\n", err)
			continue
		}

		if publishedAt.After(oneDayAgo) {
			if item.Snippet.ResourceId.VideoId == "" {
				fmt.Println("Skipping due to no video Id")
				continue
			}
			videoList[item.Snippet.ResourceId.VideoId] = &VideoInfo{item.Snippet.Title, 0}
		}

		vidResp, errGetVideo := client.Videos.List([]string{"contentDetails"}).
			Id(item.Snippet.ResourceId.VideoId).
			Do()

		if errGetVideo != nil {
			return nil, fmt.Errorf("Error get content details in video: %v", errGetVideo)
		}

		for _, video := range vidResp.Items {
			toMinutes, errConvertMinutes := iso8601ToMinutes(video.ContentDetails.Duration)

			if errConvertMinutes != nil {
				return nil, fmt.Errorf(
					"Error at convert length iso format to minutes: %v",
					errConvertMinutes,
				)
			}

			videoInfo := videoList[item.Snippet.ResourceId.VideoId]

			if videoInfo.Title == "" {
				log.Println("Skip video with no title")
				continue
			}

			if videoInfo != nil {
				videoInfo.LengthMins = toMinutes
			}

			fmt.Printf(
				"Found video title: %v, duration: %d minutes\n",
				videoInfo.Title,
				videoInfo.LengthMins,
			)
		}
	}
	fmt.Println("Exit loop")

	fmt.Printf("Video list info before export: %v\n", &videoList)

	return videoList, nil
}
