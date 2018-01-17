package main

import (
	"log"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
)

func getUserTweets(httpClient *http.Client, userToTrack string, lastTweetId int64) []TimelineTweet {
	log.Printf("Retrieving tweets for user %s starting from ID %d\n", userToTrack, lastTweetId)

	timelineUrl := fmt.Sprintf("%s%s?tweet_mode=extended&screen_name=%s&include_rts=false", TWITTER, USER_TIMELINE, userToTrack)
	if lastTweetId > 0 {
		timelineUrl = fmt.Sprintf("%s&since_id=%d", timelineUrl, lastTweetId)
	} else {
		timelineUrl += "&count=1"
	}

	res, err := httpClient.Get(timelineUrl)
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

func postTweet(httpClient *http.Client, userToTrack string, message string) {
	log.Printf("Sending tweet to user %s with message \"%s\"\n", userToTrack, message)
	statusUpdateUrl := fmt.Sprintf("%s%s", TWITTER, STATUS_UPDATE)
	res, err := httpClient.PostForm(statusUpdateUrl, url.Values{"status": {"@faesaurus testing testing 123"}, "in_reply_to_status_id": {"951931421181935616"}})
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
	if body, err := ioutil.ReadAll(res.Body); err == nil {
		if err2 := json.Unmarshal(body, &target); err2 != nil {
			return fmt.Errorf("Failed unmarshalling response: %s", err2)
		}
	
		return nil
	} else {
		return fmt.Errorf("Failed to read response body: %s", err)
	}	
}
