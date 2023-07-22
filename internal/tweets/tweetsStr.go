package tweets

import (
	"github.com/go-playground/validator/v10"
	"time"
)

type Tweet struct {
	TweetID       int       `json:"tweet_id"`
	UserID        int       `json:"user_id"`
	Author        string    `json:"author"`
	Text          string    `json:"text"`
	CreatedAt     time.Time `json:"created_at"`
	LikeCount     int       `json:"like_count"`
	Retweet       int       `json:"repost"`
	LoginToken    string    `json:"-"`
	ParentTweetId int       `json:"parent_tweet_id"`
	Visibility
}
type CreateNewTweetRequest struct {
	TweetID   int       `json:"-"`
	Text      string    `json:"text" validate:"required,checkTweetText"`
	CreatedAt time.Time `json:"-"`
	Visibility
	ReplyTo int `json:"-"`
}

type TweetValid struct {
	Validate *validator.Validate
	ValidErr map[string]string
}
