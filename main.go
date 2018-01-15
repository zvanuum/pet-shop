package main

import (
	"time"
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"
	"net/url"
	// "bytes"
	"os"
	"os/signal"
	"syscall"

	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/clientcredentials"
	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

var lastTweetId int64

func main() {
	viper.SetConfigFile("./config/config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file, %s", err)
	}

	userToTrack := viper.GetString("USER_TO_TRACK")

	// oauthConfig := &clientcredentials.Config{
    //     ClientID:     viper.GetString("CONSUMER_KEY"),
    //     ClientSecret: viper.GetString("CONSUMER_SECRET"),
    //     TokenURL:     "https://api.twitter.com/oauth2/token",
	// }
	// httpClient := oauthConfig.Client(oauth2.NoContext)
	
    config := oauth1.NewConfig(viper.GetString("CONSUMER_KEY"), viper.GetString("CONSUMER_SECRET"))
	token := oauth1.NewToken(viper.GetString("ACCESS_TOKEN"), viper.GetString("ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	ticker := time.NewTicker(time.Duration(viper.GetInt("FREQUENCY")) * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
			case <-ticker.C:
				respondToTweets(userToTrack, httpClient);
			case <-sigs:
				ticker.Stop();
				fmt.Println("Exiting")
				return;
		}
	}
}

func respondToTweets(userToTrack string, httpClient *http.Client) {
	tweets := getUserTweets(userToTrack, httpClient);
	if len(tweets) > 0 {
		lastTweetId = tweets[len(tweets) - 1].Id
	}
	
	for _, tweet := range tweets {
		fmt.Printf("\n%v\n", tweet)
	}

	postTweet(httpClient)
}


func getUserTweets(userToTrack string, httpClient *http.Client) []TimelineTweet {
	timelineUrl := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + userToTrack + "&include_rts=false"
	if lastTweetId > 0 {
		timelineUrl += fmt.Sprintf("&since_id=%d", lastTweetId)
	} else {
		timelineUrl += "&count=1"
	}

	res, err := httpClient.Get(timelineUrl)
    if err != nil {
        fmt.Printf("Failed: %s", err)
	}
	defer res.Body.Close()

	var data []TimelineTweet

	err = unmarshalResponse(res, &data)
	if err != nil {
		fmt.Printf("Couldn't get tweets: %s\n", err)
		return nil
	} 
	
	return data
}

func postTweet(httpClient *http.Client) {
	tweetUrl := "https://api.twitter.com/1.1/statuses/update.json"
	// tweet, err := json.Marshal(StatusUpdate{
	// 	Status: "@faesaurus testing testing 123",
	// 	InReplyToStatusId: int64(951931421181935616),
	// })
	// fmt.Printf("\n\n%v\n\n", string(tweet))
	// if err != nil {
	// 	fmt.Printf("Failed to create POST body, %s", err)
	// 	return
	// }

	// res, err := httpClient.Post(tweetUrl, "application/json", bytes.NewBuffer(tweet))
	res, err := httpClient.PostForm(tweetUrl, url.Values{"status": {"@faesaurus testing testing 123"}, "in_reply_to_status_id": {"951931421181935616"}})
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Failed to send tweet: %s", err)
	} else {
		bodyBytes, err2 := ioutil.ReadAll(res.Body)
		if err2 != nil {
			fmt.Printf("Failed to read response after updating status, %s", err2)
		} else {
			fmt.Printf("%d  %s    %s\n", res.StatusCode, res.Status, string(bodyBytes))
		}
	}
}

func unmarshalResponse(res *http.Response, target interface{}) error {
	if body, err := ioutil.ReadAll(res.Body); err == nil {
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return fmt.Errorf("Failed unmarshalling response: %s", err2)
		}
	
		return nil
	} else {
		return fmt.Errorf("Failed to read response body: %s", err)
	}	
}
