package admin

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"context"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

type UsersList struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

func BlockUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	err := UpdateUserBlockStatus(r.Context(), true, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, fmt.Sprintf("user %s was blocked", userID))
}

func UnblockUser(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["id_user"]
	err := UpdateUserBlockStatus(r.Context(), false, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	services.ReturnJSON(w, http.StatusOK, fmt.Sprintf("user %s was unblocked", userID))
}

func GetAllBlockUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := GetAllUsers(r.Context(), true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}

func GetAllUnblockUsers(w http.ResponseWriter, r *http.Request) {
	userList, err := GetAllUsers(r.Context(), false)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	services.ReturnJSON(w, http.StatusOK, userList)
}

func UpdateUserBlockStatus(ctx context.Context, status bool, userID string) error {
	query := "UPDATE users_tweeter SET block = $1 WHERE id = $2"
	_, err := pg.DB.ExecContext(ctx, query, status, userID)
	if err != nil {
		return err
	}
	return nil
}

func GetAllUsers(ctx context.Context, block bool) ([]UsersList, error) {
	var usersList []UsersList

	query := "SELECT id, users_tweeter FROM users_tweeter WHERE block = $1"
	rows, err := pg.DB.QueryContext(ctx, query, block)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
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
