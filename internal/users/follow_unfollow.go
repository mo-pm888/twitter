package users

import (
	"fmt"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

func (s *Service) Follow(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	secondUserID := mux.Vars(r)["id"]

	var exists bool
	err := s.DB.QueryRow("SELECT EXISTS (SELECT 1 FROM follower WHERE follower = $1 AND following = $2 LIMIT 1)", userID, secondUserID).Scan(&exists)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		services.ReturnErr(w, "You are already following this user", http.StatusBadRequest)
		return
	} else {

		_, err = s.DB.Exec("INSERT INTO follower (follower, following) VALUES ($1, $2)", userID, secondUserID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		message := fmt.Sprintf("You are now following a user with id %s", secondUserID)
		services.ReturnJSON(w, http.StatusOK, message)
	}
}

func (s *Service) Unfollow(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	secondUserID := mux.Vars(r)["id"]

	_, err := s.DB.Exec("DELETE FROM follower WHERE follower = $1 AND following = $2", userID, secondUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	message := fmt.Sprintf("You are no longer following a user with id %s", secondUserID)
	services.ReturnJSON(w, http.StatusOK, message)
}
