package tweets

import (
	"Twitter_like_application/internal/services"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

func Reply(w http.ResponseWriter, r *http.Request) {
	idTweet := mux.Vars(r)["id_tweet"]
	userID := r.Context().Value("userID").(int)
	var newReply CreatNewTweet
	err := json.NewDecoder(r.Body).Decode(&newReply)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	err = services.CreatNewTweet(newReply, r, userID, idTweet)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	massage := "Reply was created"
	services.ReturnJSON(w, http.StatusCreated, massage)

}
