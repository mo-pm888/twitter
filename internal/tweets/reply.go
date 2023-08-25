package tweets

import (
	"encoding/json"
	"net/http"
	"strconv"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type ReplyRequest struct {
	ParentID string `json:"parent_id"`
	Text     string `json:"message"`
}
type Reply struct {
	CreatNewTweet
}

func (s *Service) Reply(w http.ResponseWriter, r *http.Request) {
	var tweetReply Reply
	err := json.NewDecoder(r.Body).Decode(&tweetReply)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	tweetID := mux.Vars(r)["id_tweet"]
	parentID, err := strconv.Atoi(tweetID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	if err = s.CreateNewTweet(&tweetReply.CreatNewTweet, r.Context(), parentID); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	answer := ReplyRequest{
		ParentID: tweetID,
		Text:     "reply was created",
	}
	services.ReturnJSON(w, http.StatusCreated, answer)
}
