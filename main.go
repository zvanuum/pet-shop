package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./config/config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file, %s", err)
	}

	userToTrack := viper.GetString("USER_TO_TRACK")
	var lastTweetID int64

	config := oauth1.NewConfig(viper.GetString("CONSUMER_KEY"), viper.GetString("CONSUMER_SECRET"))
	token := oauth1.NewToken(viper.GetString("ACCESS_TOKEN"), viper.GetString("ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	ticker := time.NewTicker(time.Duration(viper.GetInt("FREQUENCY")) * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			tweets := getUserTweets(httpClient, userToTrack, lastTweetID)
			if len(tweets) == 0 {
				break
			}

			lastTweetID = tweets[0].ID

			for _, tweet := range tweets {
				log.Printf("%v\n", tweet)
				message := fmt.Sprintf("@%s %d more characters contributing to making our nation a joke. Way to go.", userToTrack, len(tweet.Text))
				postTweet(httpClient, message, tweet.ID)
			}
		case <-sigs:
			ticker.Stop()
			log.Println("Exiting")
			return
		}
	}
}
