package main

import (
	"errors"
	"fmt"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/admin"
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/server"
	"Twitter_like_application/internal/tweets"
	"Twitter_like_application/internal/users"
	"Twitter_like_application/migrations"
)

const (
	migrationGoodMsg = "migrations start"
	configBadMsg     = "reading config is mistake "
	settingsGoodMsg  = "settings read and applied "
)

func Run() error {
	c, err := config.New()
	if err != nil {
		return errors.New(configBadMsg)
	}
	db, err := pg.ConnectPostgresql(*c)
	if err != nil {
		return err
	}
	u := users.New(db)
	t := tweets.New(db)
	a := admin.New(db)
	err = migrations.Run(db)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(migrationGoodMsg)
	}
	if err = ReadSettings(a); err != nil {
		fmt.Println(err)
	}
	fmt.Println(settingsGoodMsg)
	if err = server.Server(*c, *u, *t, *a); err != nil {
		return err
	}
	return nil
}

func ReadSettings(s *admin.Service) error {
	tweetSetting, err := s.GetSettings("tweet")
	if err != nil {
		return err
	}
	s.TweetLength = tweetSetting.TweetLength
	return nil
}
