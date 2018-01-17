package main

import (
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dghubble/oauth1"
	"github.com/spf13/viper"
)

func main() {
	var spellChecker SpellChecker
	spellChecker.NewSpellChecker()

	viper.SetConfigFile("./config/config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Couldn't read configuration file, %s", err)
	}

	userToTrack := viper.GetString("USER_TO_TRACK")
	var lastTweetId int64

    config := oauth1.NewConfig(viper.GetString("CONSUMER_KEY"), viper.GetString("CONSUMER_SECRET"))
	token := oauth1.NewToken(viper.GetString("ACCESS_TOKEN"), viper.GetString("ACCESS_TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	ticker := time.NewTicker(time.Duration(viper.GetInt("FREQUENCY")) * time.Second)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
			case <-ticker.C:
				tweets := getUserTweets(httpClient, userToTrack, lastTweetId)
				if len(tweets) == 0 {
					break
				}

				lastTweetId = tweets[0].Id
				for _, tweet := range tweets {
					log.Printf("%v\n", tweet)
					misspelledCount := spellChecker.CountMisspelledWords(tweet.Text)
					log.Printf("Found %d misspelled words in tweet \"%s\"\n", misspelledCount, tweet.Text)

					// postTweet(httpClient)
				}
			case <-sigs:
				ticker.Stop();
				log.Println("Exiting")
				return;
		}
	}
}
