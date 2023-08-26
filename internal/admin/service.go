package admin

import "database/sql"

type Service struct {
	DB *sql.DB
	Settings
}
type Settings struct {
	TweetLength string `json:"max_tweet_length"`
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}
