package admin

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func BlockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	err := UpdateTweetBlockStatus(r.Context(), true, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, fmt.Sprintf("tweet %s was blocked", tweetID))
}
func UnblockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	err := UpdateTweetBlockStatus(r.Context(), false, tweetID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, fmt.Sprintf("tweet %s was unblocked", tweetID))
}

func UpdateTweetBlockStatus(ctx context.Context, status bool, tweetID string) error {
	var userID string
	query := "UPDATE tweets SET block = $1 WHERE tweet_id = $2"
	_, err := pg.DB.ExecContext(ctx, query, status, tweetID)
	if err != nil {
		return err
	}
	return nil
}
