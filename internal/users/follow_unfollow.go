package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func FollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	secondUserID := mux.Vars(r)["id"]

	var exists bool
	err := pg.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM follower WHERE follower = $1 AND following = $2 LIMIT 1)", userID, secondUserID).Scan(&exists)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		services.ReturnErr(w, "User is already following to this user", http.StatusBadRequest)
		return
	} else {

		_, err = pg.DB.Exec("INSERT INTO follower (follower, following) VALUES ($1, $2)", userID, secondUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		message := fmt.Sprintf("id %d follower to id %s", userID, secondUserID)
		services.ReturnJSON(w, http.StatusOK, message)
	}
}

func UnfollowUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	secondUserID := mux.Vars(r)["id"]

	_, err := pg.DB.Exec("DELETE FROM follower WHERE follower = $1 AND following = $2", userID, secondUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("id %d unfollower from id %s", userID, secondUserID)
	services.ReturnJSON(w, http.StatusOK, message)
}
