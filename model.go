package main

const (
	TWITTER       = "https://api.twitter.com/1.1"
	USER_TIMELINE = "/statuses/user_timeline.json"
	STATUS_UPDATE = "/statuses/update.json"
)

type TimelineTweet struct {
	CreatedAt string `json:"createdAt"`
	Id int64 `json:"id"`
	Text string `json:"full_text"`
}

type StatusUpdate struct {
	Status string `json:"status"`
	InReplyToStatusId int64 `json:"in_reply_to_status_id,omitempty"`
}
