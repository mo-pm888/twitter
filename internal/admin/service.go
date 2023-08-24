package admin

import "database/sql"

type Service struct {
	DB *sql.DB
	Settings
}
type Settings struct {
	TweetLength string `json:"length"`
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}
