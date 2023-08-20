package admin

import (
	"context"
	"fmt"
	"net/http"

	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type message struct {
	TweetID string `json:"tweet_id"`
	Message string `json:"message"`
}

func BlockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	err := UpdateTweetBlockStatus(r.Context(), true, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	m := message{
		TweetID: tweetID,
		Message: fmt.Sprintf("tweet %s was blocked", tweetID),
	}
	services.ReturnJSON(w, http.StatusOK, m)
}
func UnblockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	err := UpdateTweetBlockStatus(r.Context(), false, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	m := message{
		TweetID: tweetID,
		Message: fmt.Sprintf("tweet %s was unblocked", tweetID),
	}
	services.ReturnJSON(w, http.StatusOK, m)
}

func UpdateTweetBlockStatus(ctx context.Context, status bool, tweetID string) error {
	query := "UPDATE tweets SET block = $1 WHERE tweet_id = $2"
	_, err := pg.DB.ExecContext(ctx, query, status, tweetID)
	if err != nil {
		return err
	}
	return nil
}
