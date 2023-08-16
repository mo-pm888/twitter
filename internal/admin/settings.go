package admin

import (
	"Twitter_like_application/internal/services"
	"Twitter_like_application/internal/tweets"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func SettingTweetLength(w http.ResponseWriter, r *http.Request) {
	newLength, err := services.StrToInt(mux.Vars(r)["new_length"])
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	tweets.MaxLengthTweet = newLength
	services.ReturnJSON(w, http.StatusOK, fmt.Sprintf("a new tweet length %s", newLength))
}
