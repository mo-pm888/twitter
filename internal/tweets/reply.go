package tweets

import (
	"context"
	"encoding/json"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

const ctxKeyTweetID = "tweetID"

type Reply struct {
	CreatTweet
}
type ReplyRequest struct {
	ParentID string `json:"parent_id"`
	Text     string `json:"message"`
}

func CreateNewReply(w http.ResponseWriter, r *http.Request) {
	var tweetReply Reply
	err := json.NewDecoder(r.Body).Decode(&tweetReply)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}

	tweetID := mux.Vars(r)["id_tweet"]

	ctx := r.Context()
	ctx = context.WithValue(ctx, ctxKeyTweetID, tweetID)

	newRequest := r.WithContext(ctx)
	err = tweetReply.Create(tweetReply.CreatTweet, newRequest.Context())
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	answer := ReplyRequest{
		ParentID: tweetID,
		Text:     "reply was created",
	}
	services.ReturnJSON(w, http.StatusCreated, answer)
}
