package users

import (
	"context"
	"database/sql"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type FollowerList struct {
	ID int `json:"id"`
}

func (s *Service) GetAllFollowers(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	userList, err := GetSubscribers(r.Context(), userID, 1, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}
func (s *Service) GetAllFollowings(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	userList, err := GetSubscribers(r.Context(), userID, 0, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}

func GetSubscribers(ctx context.Context, userID string, status int, s *sql.DB) ([]FollowerList, error) {
	var (
		usersList []FollowerList
		query     string
	)
	if status == 1 {
		query = "SELECT following FROM follower WHERE follower = $1"
	} else {
		query = "SELECT follower FROM follower WHERE following = $1"
	}
	rows, err := s.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var user FollowerList
		err = rows.Scan(&user.ID)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}
