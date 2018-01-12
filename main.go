package main

import (
	"time"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./config/config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file, %s", err)
	}

	userToTrack := viper.GetString("USER_TO_TRACK")

	oauthConfig := &clientcredentials.Config{
        ClientID:     viper.GetString("CONSUMER_KEY"),
        ClientSecret: viper.GetString("CONSUMER_SECRET"),
        TokenURL:     "https://api.twitter.com/oauth2/token",
	}
	httpClient := oauthConfig.Client(oauth2.NoContext)

	ticker := time.NewTicker(time.Duration(viper.GetInt("FREQUENCY")) * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
			case <-ticker.C:
				getUserTweets(userToTrack, httpClient);
			case <-sigs:
				ticker.Stop();
				fmt.Println("Exiting")
				return;
		}
	}
}

func unmarshalResponse(res *http.Response) ([]TimelineTweet, error) {
	if body, err := ioutil.ReadAll(res.Body); err == nil {
		var target []TimelineTweet
	
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return nil, fmt.Errorf("Failed unmarshalling response for user tweets: %s", err2)
		}
	
		return target, nil
	} else {
		return nil, fmt.Errorf("Failed to read response body: %s", err)
	}	
}

func getUserTweets(userToTrack string, httpClient *http.Client) {
	timelineUrl := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + userToTrack + "&count=1&include_rts=false"
	resp, err := httpClient.Get(timelineUrl)
    if err != nil {
        fmt.Printf("Failed: %s", err)
	}

	data, err := unmarshalResponse(resp)
	if err != nil {
		fmt.Printf("Couldn't get tweets: %s\n", err)
	} else {	
		fmt.Printf("Response: %v\n", data)
	}
}
