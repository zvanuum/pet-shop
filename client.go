package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func getUserTweets(httpClient *http.Client, userToTrack string, lastTweetID string) []TimelineTweet {
	log.Printf("Retrieving tweets for user %s starting from ID %s\n", userToTrack, lastTweetID)

	timelineURL := fmt.Sprintf("%s%s?tweet_mode=extended&screen_name=%s&include_rts=false", Twitter, UserTimeline, userToTrack)
	if lastTweetID != "" {
		timelineURL = fmt.Sprintf("%s&since_id=%s", timelineURL, lastTweetID)
	} else {
		timelineURL += "&count=1"
	}

	res, err := httpClient.Get(timelineURL)
	if err != nil {
		log.Printf("Failed to retrieve user timeline: %s", err)
	}
	defer res.Body.Close()

	var data []TimelineTweet

	err = unmarshalResponse(res, &data)
	if err != nil {
		log.Printf("Failed to unmarshal user timeline tweets: %s\n", err)
		return nil
	}

	return data
}

func postTweet(httpClient *http.Client, message string, inResponseToID string) {
	log.Printf("Sending tweet with message \"%s\" in response to status %s\n", message, inResponseToID)

	statusUpdateURL := fmt.Sprintf("%s%s", Twitter, StatusUpdate)
	res, err := httpClient.PostForm(statusUpdateURL, url.Values{"status": {message}, "in_reply_to_status_id": {inResponseToID}})
	if err != nil {
		log.Printf("Failed to send tweet: %s", err)
		return
	}

	defer res.Body.Close()
	bodyBytes, err2 := ioutil.ReadAll(res.Body)
	if err2 != nil {
		log.Printf("Failed to read response after updating status, %s", err2)
	} else {
		log.Printf("%d  %s    %s\n", res.StatusCode, res.Status, string(bodyBytes))
	}
}

func unmarshalResponse(res *http.Response, target interface{}) error {
	var body []byte
	var err error

	if body, err = ioutil.ReadAll(res.Body); err == nil {
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return fmt.Errorf("Failed unmarshalling response: %s", err2)
		}

		return nil
	}

	return fmt.Errorf("Failed to read response body: %s", err)
}
