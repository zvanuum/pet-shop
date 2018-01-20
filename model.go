package main

const (
	//Twitter API URL
	Twitter = "https://api.twitter.com/1.1"
	//UserTimeline resource for GETting a user's tweets
	UserTimeline = "/statuses/user_timeline.json"
	//StatusUpdate resource for POSTing a tweet
	StatusUpdate = "/statuses/update.json"
)

// TimelineTweet Struct to represent a tweet retrieved from a user's timeline
type TimelineTweet struct {
	CreatedAt string `json:"createdAt"`
	ID        int64  `json:"id"`
	Text      string `json:"full_text"`
}

// StatusUpdateBody Struct to represent the body of a tweet
type StatusUpdateBody struct {
	Status            string `json:"status"`
	InReplyToStatusID int64  `json:"in_reply_to_status_id,omitempty"`
}
