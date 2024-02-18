package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)

	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to %s\n", path)

	f, errOpenFile := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)

	if errOpenFile != nil {
		log.Fatalf("Unable to cache auth token: %v", errOpenFile)
	}

	defer f.Close()

	json.NewEncoder(f).Encode(token)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)

	fmt.Printf("Go to the following link in your browser then pass the auth code:\n %v \n", authURL)

	var authCode string

	if _, errReadAuthCode := fmt.Scan(&authCode); errReadAuthCode != nil {
		log.Fatalf("Unable to read auth code: %v", errReadAuthCode)
	}

	tok, errGetToken := config.Exchange(context.Background(), authCode)

	if errGetToken != nil {
		log.Fatalf("Cannot get token: %v", errGetToken)
	}

	return tok
}

func getClient(config *oauth2.Config) *http.Client {
	tokFile := "token.json"

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}

	return config.Client(context.Background(), tok)
}

func GetAuth() (*youtube.Service, error) {
	ctx := context.Background()

	jsonCred, errReadCredentials := os.ReadFile("credentials.json")

	if errReadCredentials != nil {
		return nil, errReadCredentials
	}

	config, err := google.ConfigFromJSON(jsonCred, youtube.YoutubeReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := getClient(config)

	service, errCreateService := youtube.NewService(ctx, option.WithHTTPClient(client))

	if errCreateService != nil {
		return nil, errCreateService
	}

	return service, nil
}
