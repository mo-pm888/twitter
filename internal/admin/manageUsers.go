package admin

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"Twitter_like_application/internal/services"

	"github.com/gorilla/mux"
)

type UsersList struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type messageRequest struct {
	UserID string `json:"id"`
	Text   string `json:"message"`
}

func (s *Service) BlockUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	err := UpdateUserBlockStatus(r.Context(), true, userID, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	message := messageRequest{
		UserID: userID,
		Text:   fmt.Sprintf("user %s was blocked", userID),
	}
	services.ReturnJSON(w, http.StatusOK, message)
}

func (s *Service) UnblockUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	err := UpdateUserBlockStatus(r.Context(), false, userID, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	message := messageRequest{
		UserID: userID,
		Text:   fmt.Sprintf("user %s was unblocked", userID),
	}
	services.ReturnJSON(w, http.StatusOK, message)
}

func (s *Service) GetAllBlockUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := GetAllUsers(r.Context(), true, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}

func (s *Service) GetAllUnblockUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := GetAllUsers(r.Context(), false, s.DB)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}

func UpdateUserBlockStatus(ctx context.Context, status bool, userID string, s *sql.DB) error {
	query := "UPDATE users_tweeter SET block = $1 WHERE id = $2"
	_, err := s.ExecContext(ctx, query, status, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers(ctx context.Context, block bool, s *sql.DB) ([]UsersList, error) {
	var usersList []UsersList

	query := "SELECT id,name FROM users_tweeter WHERE block = $1"
	rows, err := s.QueryContext(ctx, query, block)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	for rows.Next() {
		var user UsersList
		err = rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return nil, err
		}
		usersList = append(usersList, user)
	}

	return usersList, nil
}
