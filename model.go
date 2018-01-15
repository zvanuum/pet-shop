package main

type TimelineTweet struct {
	CreatedAt string `json:"createdAt"`
	Id int64 `json:"id"`
	Text string `json:"text"`
}

type StatusUpdate struct {
	Status string `json:"status"`
	InReplyToStatusId int64 `json:"in_reply_to_status_id,omitempty"`
}
