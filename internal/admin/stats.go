package admin

import (
	"net/http"

	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
)

type counts struct {
	countTweets int
	countUsers  int
}

func Stats(w http.ResponseWriter, r *http.Request) {
	if r.Context().Value("isAdmin").(bool) == true {
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
	} else {
		services.ReturnJSON(w, http.StatusUnauthorized, "You aren't an administrator")
	}
}
