package tweets

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

type CreateNewTweet struct {
	TweetID             int
	Text                string `json:"text" validate:"required,checkTweetText"`
	CreatedAt           time.Time
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
	Visibility
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	var newTweet CreateNewTweet
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

func (s *Service) CreateNewTweet(tweet *CreateNewTweet, ctx context.Context, parentID int) error {
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
	userID := ctx.Value("userID").(string)
	id, err := strconv.Atoi(userID)
	if err != nil {
		return err
	}
	if !tweet.isValid() {
		return errors.New("visibility error")
	}
	query := `INSERT INTO tweets (user_id, text, created_at,parent_tweet_id, public, only_followers, only_mutual_followers, only_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8) RETURNING tweet_id`
	err = s.DB.QueryRowContext(ctx, query, id, tweet.Text, time.Now(), parentID, tweet.Public, tweet.OnlyFollowers, tweet.OnlyMutualFollowers, tweet.OnlyMe).Scan(&tweet.TweetID)
	if err != nil {
		return err
	}
	return nil

}
