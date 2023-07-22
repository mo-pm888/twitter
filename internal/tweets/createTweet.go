package tweets

import (
	"Twitter_like_application/internal/services"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func CreateNewTweet(w http.ResponseWriter, r *http.Request) {
	tweetValid := &TweetValid{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := RegisterTweetValidations(tweetValid); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	userID := r.Context().Value("userID").(int)
	var newTweet CreateNewTweetRequest
	err := json.NewDecoder(r.Body).Decode(&newTweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = CreateTweet(newTweet, r.Context(), userID, "")
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTweet)

	return
}
