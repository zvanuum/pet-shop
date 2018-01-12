package main

import (
	"encoding/json"
	"fmt"
	"log"
	"io/ioutil"
	"net/http"

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
	token, err := oauthConfig.Token(oauth2.NoContext)
	httpClient := oauthConfig.Client(oauth2.NoContext)
	
	timelineUrl := "https://api.twitter.com/1.1/statuses/user_timeline.json?screen_name=" + userToTrack + "&count=3&include_rts=false"
	resp, err := httpClient.Get(timelineUrl)
    if (err != nil) {
        fmt.Printf("Failed: %s", err)
	}

	data, err := unmarshalResponse(resp)
	if err != nil {
		fmt.Printf("Couldn't get tweets: %s\n", err)
	} else {	
		fmt.Printf("Token: %s\n\n\nResponse: %v\n", token, data)
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

// func unmarshalResponse(res *http.Response) ([]map[string]interface{}, error) {
// 	if body, err := ioutil.ReadAll(res.Body); err == nil {
// 		var testing []map[string]interface{}
	
// 		if err2 := json.Unmarshal(body, &testing); err2 != nil {
// 			return nil, fmt.Errorf("Failed unmarshalling response for user tweets: %s", err2)
// 		}
	
// 		return testing, nil
// 	} else {
// 		return nil, fmt.Errorf("Failed to read response body: %s", err)
// 	}
// }

// func getJson(url string, target interface{}) error {
//     r, err := myClient.Get(url)
//     if err != nil {
//         return err
//     }
//     defer r.Body.Close()

//     return json.NewDecoder(r.Body).Decode(target)
// }
