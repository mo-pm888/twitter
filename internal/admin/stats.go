package admin

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"net/http"
)

type counts struct {
	countTweets int
	countUsers  int
}

func Stats(w http.ResponseWriter, r *http.Request) {
	var allCounts counts
	err := pg.DB.QueryRow("SELECT COUNT(*) FROM users_tweeter").Scan(&allCounts.countUsers)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}

	err = pg.DB.QueryRow("SELECT COUNT(*) FROM tweets").Scan(&allCounts.countTweets)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, allCounts)
}
