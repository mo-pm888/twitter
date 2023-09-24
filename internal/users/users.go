package users

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"net/url"

	"Twitter_like_application/internal/services"
)

const (
	ctxKeyUserID  = "userID"
	ctxKeyIsAdmin = "isAdmin"
)

func checkAuth(w http.ResponseWriter, r *http.Request, s *sql.DB) *http.Request {
	apikey := r.Header.Get("X-API-KEY")
	cookie, err := r.Cookie("session")
	if apikey == "" && (err != nil || cookie == nil) {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	var sessionID string
	if apikey != "" {
		sessionID = apikey
	} else if cookie != nil {
		sessionID = cookie.Value
	}
	if sessionID == "" {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	query := `SELECT us.user_id, ut.admin FROM user_session us JOIN users_tweeter ut ON us.user_id = ut.id WHERE us.session_id = $1 LIMIT 1`

	var userID int
	var isAdmin bool
	err = s.QueryRow(query, sessionID).Scan(&userID, &isAdmin)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyIsAdmin, isAdmin)
	return r.WithContext(ctx)
}

func (s *Service) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = checkAuth(w, r, s.DB)
		if v, ok := r.Context().Value(ctxKeyUserID).(string); ok && v != "" {
			next.ServeHTTP(w, r)
		} else {
			services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		}
	})

}

func (s *Service) AdminAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = checkAuth(w, r, s.DB)
		if v, ok := r.Context().Value(ctxKeyIsAdmin).(bool); ok && v {
			next.ServeHTTP(w, r)
		} else {
			services.ReturnErr(w, "Unauthorized as an admin", http.StatusUnauthorized)
		}
	})
}

func (s *Service) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var userResPass User
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "SELECT name, email  FROM users_tweeter WHERE id = $1"
	var user User
	err = s.DB.QueryRow(query, userID).Scan(&user.Name, &user.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userResPass.Email = user.Email
	userResPass.Name = user.Name

	ResetPasswordPlusEmail(&userResPass)
}

func ResetPasswordPlusEmail(user *User) {

	resetToken := services.GenerateResetToken()
	user.ResetPasswordToken = resetToken
	confirmURL := &url.URL{
		Scheme: "http",
		Host:   "test.com",
		Path:   "/reset-password",
		RawQuery: url.Values{
			"token": {resetToken},
		}.Encode(),
	}
	to := user.Email
	subject := "Reset your password"
	body := fmt.Sprintf("Dear %s,\n\nReset your password: click this link:\n%s", user.Name, confirmURL.String())

	auth := smtp.PlainAuth("", "your-email", "password", "your-site")
	err := smtp.SendMail("your-site:587", auth, "your-email", []string{to}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return
	}
	return
}

func (s *Service) SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, name, username FROM users WHERE name ILIKE $1 OR username ILIKE $1"

	rows, err := s.DB.Query(query, searchQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func EmailVerificationToken(email string) string {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return ""
	}
	confirmToken := base64.URLEncoding.EncodeToString(token)

	confirmURL := &url.URL{
		Scheme: "http",
		Host:   "test.com",
		Path:   "/confirm-email",
		RawQuery: url.Values{
			"token": {confirmToken},
		}.Encode(),
	}
	subject := "Confirment your email"
	body := fmt.Sprintf("Confirment email: click this link:\n%s", confirmURL.String())

	auth := smtp.PlainAuth("", "your email", "password", "your site/token")

	err = smtp.SendMail("your email:587", auth, "your site/token", []string{email}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return ""
	}

	return confirmToken
}

func DeleteUserSession(token string, s *sql.DB) error {
	query := "DELETE FROM user_session WHERE login_token = $1"
	_, err := s.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
