package tweets

import (
	"database/sql"
	"time"

	"github.com/go-playground/validator/v10"
)

const (
	defaultPublic              = true
	defaultOnlyFollowers       = false
	defaultOnlyMutualFollowers = false
	defaultOnlyMe              = false
)

type Service struct {
	DB *sql.DB
}
type Visibility struct {
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}
type Tweet struct {
	TweetID       int       `json:"tweet_id"`
	UserID        int       `json:"user_id"`
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

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}

func (v *Visibility) count() int {
	count := 0
	switch true {
	case v.Public:
		count++
	case v.OnlyFollowers:
		count++
	case v.OnlyMutualFollowers:
		count++
	case v.OnlyMe:
		count++
	}
	return count
}
func (v *Visibility) isValid() bool {
	return v.count() < 2
}
func (v *Visibility) defaultVisibilities() bool {
	return v.count() == 0
}

//func (v *Visibility) defaultVisibilities(tweet *Tweet) {
//	if tweet.Visibility.Public == nil {
//		tweet.Public = &defaultPublic
//	}
//	if tweet.OnlyFollowers == nil {
//		tweet.OnlyFollowers = &defaultOnlyFollowers
//	}
//	if tweet.OnlyMutualFollowers == nil {
//		tweet.OnlyMutualFollowers = &defaultOnlyMutualFollowers
//	}
//	if tweet.OnlyMe == nil {
//		tweet.OnlyMe = &defaultOnlyMe
//	}
//}
