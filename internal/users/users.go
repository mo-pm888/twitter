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
	"strings"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

const (
	ctxKeyUserID  = "userID"
	ctxKeyIsAdmin = "isAdmin"
)

type UserValid struct {
	validate *validator.Validate
	validErr map[string]string
}
type NameVal struct {
	short    bool
	long     bool
	realName bool
}

func (v *UserValid) Error() string {
	var pairs []string
	for k, v := range v.validErr {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(pairs, "; ")
}

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

	var userID string
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

//func GetFollowers(w http.ResponseWriter, r *http.Request) {
//	userID := r.FormValue("user_id")
//	if userID == "" {
//		http.Error(w, "Missing user ID", http.StatusBadRequest)
//		return
//	}
//
//	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.follower_id WHERE s.followee_id = $1"
//	rows, err := pg.DB.Query(query, userID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer rows.Close()
//
//	var followers []User
//
//	for rows.Next() {
//		var follower User
//		err := rows.Scan(&follower.UserID, &follower.Name)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		followers = append(followers, follower)
//	}
//
//	if err := rows.Err(); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(followers)
//}

//func GetFollowing(w http.ResponseWriter, r *http.Request) {
//	userID := r.FormValue("user_id")
//	if userID == "" {
//		http.Error(w, "Missing user ID", http.StatusBadRequest)
//		return
//	}
//
//	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.followee_id WHERE s.follower_id = $1"
//	rows, err := pg.DB.Query(query, userID)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	defer rows.Close()
//
//	var following []User
//
//	for rows.Next() {
//		var followee User
//		err := rows.Scan(&followee.UserID, &followee.Name)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//		following = append(following, followee)
//	}
//
//	if err := rows.Err(); err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(following)
//}

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

//	func GetStatistics(w http.ResponseWriter, r *http.Request) {
//		userCount, err := services.GetUserCount()
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		tweetCount, err := services.GetTweetCount()
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		statistics := struct {
//			UserCount  int `json:"user_count"`
//			TweetCount int `json:"tweet_count"`
//		}{
//			UserCount:  userCount,
//			TweetCount: tweetCount,
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(statistics)
//	}
func CheckEmail(newUser *User) string {
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
	to := newUser.Email
	subject := "Confirment your email"
	body := fmt.Sprintf("Confirment email: click this link:\n%s", confirmURL.String())

	auth := smtp.PlainAuth("", "your email", "password", "your site/token")

	err = smtp.SendMail("your email:587", auth, "your site/token", []string{to}, []byte(fmt.Sprintf("Subject: %s\n\n%s", subject, body)))
	if err != nil {
		return ""
	}

	return confirmToken
}

//	func GetUserProfile(w http.ResponseWriter, r *http.Request) {
//		vars := mux.Vars(r)
//		userID := vars["id"]
//
//		query := "SELECT id, name, bio FROM users_tweeter WHERE id = $1"
//		var user User
//		err := pg.DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Bio)
//		if err != nil {
//			http.Error(w, err.Error(), http.StatusInternalServerError)
//			return
//		}
//
//		response := struct {
//			Name     string `json:"name"`
//			Email    string `json:"email"`
//			Birthday string `json:"birthday"`
//			NickName string `json:"nickName"`
//			Bio      string `json:"bio"`
//			Location string `json:"location"`
//		}{
//			Name:     user.Name,
//			Email:    user.Email,
//			Birthday: user.BirthDate,
//			NickName: user.Nickname,
//			Bio:      user.Bio,
//			Location: user.Location,
//		}
//
//		w.Header().Set("Content-Type", "application/json")
//		json.NewEncoder(w).Encode(response)
//	}
func DeleteUserSession(token string, s *sql.DB) error {
	query := "DELETE FROM user_session WHERE login_token = $1"
	_, err := s.Exec(query, token)
	if err != nil {
		return err
	}

	return nil
}
func RegisterUsersValidations(userValid *UserValid) error {
	err := userValid.validate.RegisterValidation("checkPassword", func(fl validator.FieldLevel) bool {
		return CheckPassword(fl, userValid)
	})
	if err != nil {
		return err
	}
	err = userValid.validate.RegisterValidation("checkName", func(fl validator.FieldLevel) bool {
		return CheckName(fl, userValid)
	})
	if err != nil {
		return err
	}

	err = userValid.validate.RegisterValidation("checkDate", func(fl validator.FieldLevel) bool {
		return CheckDate(fl, userValid)
	})
	if err != nil {
		return err
	}
	err = userValid.validate.RegisterValidation("checkNickname", func(fl validator.FieldLevel) bool {
		return CheckNickName(fl, userValid)
	})
	if err != nil {
		return err
	}
	err = userValid.validate.RegisterValidation("checkBio", func(fl validator.FieldLevel) bool {
		return CheckBio(fl, userValid)
	})
	if err != nil {
		return err
	}
	err = userValid.validate.RegisterValidation("checkLocation", func(fl validator.FieldLevel) bool {
		return CheckLocation(fl, userValid)
	})
	if err != nil {
		return err
	}
	err = userValid.validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return CheckEmailVal(fl, userValid)
	})
	if err != nil {
		return err
	}
	return nil
}
