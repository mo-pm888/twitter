package admin

import (
	"net/http"

	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
)

type counts struct {
	CountTweets int `json:"tweets"`
	CountUsers  int `json:"users"`
}

func Stats(w http.ResponseWriter, r *http.Request) {
	var allCounts counts
	err := pg.DB.QueryRow("SELECT COUNT(*) FROM users_tweeter").Scan(&allCounts.CountUsers)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}

	err = pg.DB.QueryRow("SELECT COUNT(*) FROM tweets").Scan(&allCounts.CountTweets)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, allCounts)
}
