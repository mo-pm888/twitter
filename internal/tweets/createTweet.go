package tweets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

type CreatTweet struct {
	TweetID   int
	Text      string `json:"text" validate:"required,checkTweetText"`
	CreatedAt time.Time
	ParentID  int
	Visibility
}
type Visibility struct {
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}

func (v *Visibility) count() int {
	count := 0
	if v.Public {
		count++
	}
	if v.OnlyFollowers {
		count++
	}
	if v.OnlyMutualFollowers {
		count++
	}
	if v.OnlyMe {
		count++
	}
	return count
}
func (v *Visibility) isValid() bool {
	return v.count() < 2
}

func (c CreatTweet) Create(tweet CreatTweet, ctx context.Context) error {
	userID := ctx.Value("userID").(string)
	tweetValid := &TweetValid{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := RegisterTweetValidations(tweetValid); err != nil {
		return err
	}
	err := tweetValid.Validate.Struct(tweet)
	if err != nil {
		return err
	}
	if !c.isValid() {
		return errors.New("there must be only one visibility parameter")
	}
	parentID, ok := ctx.Value("tweetID").(string)
	if !ok {
		parentID = "0"
	}
	parentIDint, err := strconv.Atoi(parentID)
	if err != nil {
		return err
	}
	query := `INSERT INTO tweets (user_id, text, created_at, public, only_followers, only_mutual_followers, only_me,parent_tweet_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8) RETURNING tweet_id`
	err = pg.DB.QueryRowContext(ctx, query, userID, tweet.Text, time.Now(), tweet.Public, tweet.OnlyFollowers, tweet.OnlyMutualFollowers, tweet.OnlyMe, parentIDint).Scan(&tweet.TweetID)
	if err != nil {
		return err

	}
	return nil

}

func CreateNewTweet(w http.ResponseWriter, r *http.Request) {
	var tweet CreatTweet
	err := json.NewDecoder(r.Body).Decode(&tweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = tweet.Create(tweet, r.Context()); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusCreated, tweet)
}
