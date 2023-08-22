package tweets

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type Tweet struct {
	TweetID       int       `json:"tweet_id"`
	UserID        int       `json:"user_id"`
	Author        string    `json:"author"`
	Text          string    `json:"text"`
	CreatedAt     time.Time `json:"created_at"`
	LikeCount     int       `json:"like_count"`
	Retweet       int       `json:"repost"`
	ParentTweetId int       `json:"parent_tweet_id"`
	Visibility
}

type TweetValid struct {
	Validate *validator.Validate
	ValidErr map[string]string
}
type ReplayTweet struct {
	Tweet
}
