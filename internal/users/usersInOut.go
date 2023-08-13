package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

const (
	PasswordMaxLength = 320
	EmailMaxLength    = 320
)

func LogIn(w http.ResponseWriter, r *http.Request) {
	request := &LoginRequest{}
	err := json.NewDecoder(r.Body).Decode(request)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if len(request.Password) > PasswordMaxLength {
		services.ReturnErr(w, "The password is too long", http.StatusBadRequest)
		return
	}
	if len(request.Email) > EmailMaxLength {
		services.ReturnErr(w, "The email is too long", http.StatusBadRequest)
		return
	}
	query := "SELECT id, password FROM users_tweeter WHERE email = $1"
	var userID int
	var savedPassword string
	err = pg.DB.QueryRow(query, request.Email).Scan(&userID, &savedPassword)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(savedPassword), []byte(request.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			services.ReturnJSON(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return

	}
	sessionID := uuid.New().String()

	cookie := &http.Cookie{
		Name:     "session",
		Value:    sessionID,
		Expires:  time.Now().AddDate(0, 0, 30),
		HttpOnly: true,
		Path:     "/",
	}
	http.SetCookie(w, cookie)

	insertQuery := "INSERT INTO user_session (user_id, login_token, timestamp) VALUES ($1, $2, $3)"
	_, err = pg.DB.Exec(insertQuery, userID, cookie.Value, time.Now())
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	services.ReturnJSON(w, http.StatusOK, "Authentication succeeded")
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	apikey := r.Header.Get("X-API-KEY")
	if apikey == "" {
		cookie, err := r.Cookie("session")
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = DeleteUserSession(cookie.Value)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusBadRequest)
			return
		}

		cookie = &http.Cookie{
			Name:    "session",
			Value:   "",
			Expires: time.Now().AddDate(0, 0, -1),
			Path:    "/",
		}
		http.SetCookie(w, cookie)
		services.ReturnJSON(w, http.StatusOK, "Exit completed successfully")
	} else {
		err := DeleteUserSession(r.Header.Get("X-API-KEY"))
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusBadRequest)
			return
		}
		services.ReturnJSON(w, http.StatusOK, "Exit completed successfully")
	}
}
