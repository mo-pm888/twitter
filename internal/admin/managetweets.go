package admin

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type moderateTweetResponse struct {
	TweetID string `json:"tweet_id"`
	Message string `json:"message"`
}

func (s *Service) BlockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	if err := s.UpdateTweetBlockStatus(r.Context(), true, tweetID); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	msg := moderateTweetResponse{
		TweetID: tweetID,
		Message: fmt.Sprintf("tweet %s was blocked", tweetID),
	}
	services.ReturnJSON(w, http.StatusOK, msg)
}
func (s *Service) UnblockTweet(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	err := s.UpdateTweetBlockStatus(r.Context(), false, tweetID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	msg := moderateTweetResponse{
		TweetID: tweetID,
		Message: fmt.Sprintf("tweet %s was unblocked", tweetID),
	}
	services.ReturnJSON(w, http.StatusOK, msg)
}

func (s *Service) UpdateTweetBlockStatus(ctx context.Context, status bool, tweetID string) error {
	query := "UPDATE tweets SET block = $1 WHERE tweet_id = $2 AND block !=$1"
	result, err := s.DB.ExecContext(ctx, query, status, tweetID)
	if err != nil {
		return err
	}
	f, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if f == 0 {
		return errors.New("tweet already has this status ")
	}
	return nil
}
