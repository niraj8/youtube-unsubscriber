package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/youtube/v3"
)

const missingClientSecretsMessage = `Please configure OAuth 2.0`

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("youtube-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func handleError(err error, message string) {
	if message == "" {
		message = "Error making API call"
	}
	if err != nil {
		log.Fatalf(message+": %v", err.Error())
	}
}

func main() {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/youtube-go-quickstart.json
	config, err := google.ConfigFromJSON(b, youtube.YoutubeReadonlyScope, youtube.YoutubeForceSslScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)
	service, err := youtube.New(client)

	handleError(err, "Error creating YouTube client")

	allSubscriptions := getAllSubscriptions(service, ctx)
	writeSubscriptionsToFile(allSubscriptions)

	// allSubscriptions, err := readSubscriptionsFromFile()
	// handleError(err, "Error reading subscriptions from file")
	fmt.Printf("Found %d subscriptions\n", len(allSubscriptions))
	unsubscribeFromAllChannels(service, ctx, allSubscriptions)
}

func unsubscribeFromAllChannels(service *youtube.Service, ctx context.Context, allSubscriptions []*youtube.Subscription) {
	for _, subscription := range allSubscriptions {
		fmt.Printf("Unsubscribing from %s\n", subscription.Snippet.Title)
		err := service.Subscriptions.Delete(subscription.Id).Do()
		handleError(err, "Error unsubscribing from channel "+subscription.Snippet.Title)
	}
}

func readSubscriptionsFromFile() ([]*youtube.Subscription, error) {
	f, err := os.Open("subscriptions.json")
	if err != nil {
		log.Fatalf("Unable to read file: %v", err)
		return nil, err
	}
	defer f.Close()

	var allSubscriptions []*youtube.Subscription
	json.NewDecoder(f).Decode(&allSubscriptions)
	return allSubscriptions, nil
}

func writeSubscriptionsToFile(allSubscriptions []*youtube.Subscription) {
	b_subscriptions, err := json.Marshal(allSubscriptions)
	if err != nil {
		log.Fatalf("Unable to convert data to json: %v", err)
		return
	}

	f, err := os.Create("subscriptions.json")
	if err != nil {
		log.Fatalf("Unable to create file: %v", err)
		return
	}
	fmt.Fprint(f, string(b_subscriptions))
	defer f.Close()
}

func getAllSubscriptions(service *youtube.Service, ctx context.Context) []*youtube.Subscription {
	allSubscriptions := make([]*youtube.Subscription, 0)
	service.Subscriptions.List([]string{"snippet,contentDetails"}).Mine(true).Pages(ctx, func(page *youtube.SubscriptionListResponse) error {
		allSubscriptions = append(allSubscriptions, page.Items...)
		return nil
	})
	return allSubscriptions
}
