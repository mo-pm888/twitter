package tweets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	var newTweet Tweet
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = s.CreateNewTweet(&newTweet, r.Context(), 0); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusCreated, newTweet)
}

func (s *Service) CreateNewTweet(tweet *Tweet, ctx context.Context, parentID int) error {
	tweetValid := &TweetValid{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := RegisterTweetValidations(tweetValid); err != nil {
		fmt.Println(err)
	}
	if err := tweetValid.Validate.Struct(tweet); err != nil {
		return err
	}
	userID := ctx.Value("userID").(int)
	if !tweet.isValid() {
		return errors.New("visibility error, many visual arguments")
	}
	if tweet.defaultVisibilities() {
		tweet.Visibility.Public = true
		tweet.Visibility.OnlyMe = false
		tweet.Visibility.OnlyFollowers = false
		tweet.Visibility.OnlyMutualFollowers = false
	}
	query := `INSERT INTO tweets (user_id, text, created_at,parent_tweet_id, public, only_followers, only_mutual_followers, only_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8) RETURNING tweet_id`
	if err := s.DB.QueryRowContext(ctx, query, userID, tweet.Text, time.Now(), parentID, tweet.Public, tweet.OnlyFollowers, tweet.OnlyMutualFollowers, tweet.OnlyMe).Scan(&tweet.TweetID); err != nil {
		return err
	}
	return nil

}
