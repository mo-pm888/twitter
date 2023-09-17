package tweets

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type editTweetRequest = createTweetRequest

type Visibility struct {
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}

func (v *Visibility) count() int {
	count := 0
	switch true {
	case v.Public:
		count++
	case v.OnlyFollowers:
		count++
	case v.OnlyMutualFollowers:
		count++
	case v.OnlyMe:
		count++
	}
	return count
}
func (v *Visibility) isValid() bool {
	return v.count() < 2
}

func (s *Service) Edit(w http.ResponseWriter, r *http.Request) {
	tweetID := mux.Vars(r)["id_tweet"]
	userID := r.Context().Value("userID").(int)
	var req editTweetRequest
	var tweet Tweet
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		services.ReturnErr(w, err, http.StatusBadRequest)
		return
	}
	if !req.isValid() {
		services.ReturnErr(w, "There must be only one visibility parameter", http.StatusInternalServerError)
		return
	}
	query := "SELECT user_id, public, only_followers, only_mutual_followers, only_me FROM tweets WHERE tweet_id = $1"
	err = s.DB.QueryRow(query, tweetID).Scan(&tweet.UserID, &tweet.Public, &tweet.OnlyFollowers, &tweet.OnlyMutualFollowers, &tweet.OnlyMe)
	if err != nil {
		if err == sql.ErrNoRows {
			services.ReturnErr(w, "Tweet not found", http.StatusNotFound)
		} else {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	if tweet.UserID != userID {
		http.Error(w, "it isn't your tweet", http.StatusUnauthorized)
		return
	}
	var visibility Visibility
	if req.count() == 0 {
		visibility = tweet.Visibility
	} else {
		visibility = req.Visibility
	}
	query = "UPDATE tweets SET text = $1, public = $2, only_followers = $3, only_mutual_followers = $4, only_me = $5 WHERE tweet_id = $6"
	_, err = s.DB.ExecContext(r.Context(), query, req.Text, visibility.Public, visibility.OnlyFollowers, visibility.OnlyMutualFollowers, visibility.OnlyMe, tweetID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response := map[string]interface{}{
		"status":  "success",
		"message": "Tweet updated successfully",
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
