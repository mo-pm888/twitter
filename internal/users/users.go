package users

import (
	_ "Twitter_like_application/internal/database/pg"
	pg "Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/smtp"
	"net/url"
	"strings"
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

func handleAuthenticatedRequest(w http.ResponseWriter, r *http.Request, next http.Handler) {
	apikey := r.Header.Get("X-API-KEY")
	cookie, err := r.Cookie("session")
	if apikey == "" && (err != nil || cookie == nil) {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var sessionID string
	if apikey != "" {
		sessionID = apikey
	} else if cookie != nil {
		sessionID = cookie.Value
	}
	if cookie != nil || apikey != "" {
		query := "SELECT user_id FROM user_session WHERE login_token = $1"
		var userID int
		err = pg.DB.QueryRow(query, sessionID).Scan(&userID)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", userID)
		r = r.WithContext(ctx)

	} else {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	next.ServeHTTP(w, r)
}

func AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleAuthenticatedRequest(w, r, next)
	})
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	newUser := &User{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `SELECT id FROM users_tweeter WHERE email = $1`
	var existingUserID int
	err = pg.DB.QueryRow(query, newUser.Email).Scan(&existingUserID)
	if err == nil {
		services.ReturnErr(w, "User with this email already exists", http.StatusBadRequest)
		return
	} else if err != sql.ErrNoRows {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.BirthDate == "" {
		services.ReturnErr(w, "Invalid user data", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	query = `INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = pg.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			services.ReturnErr(w, "This user is already added", http.StatusBadRequest)
			return
		}
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userToken := CheckEmail(newUser)
	newUser.EmailToken = userToken

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newUser)
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)

	var userResPass User
	err := json.NewDecoder(r.Body).Decode(&userResPass)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := "SELECT name, email  FROM users_tweeter WHERE id = $1"
	var user User
	err = pg.DB.QueryRow(query, userID).Scan(&user.Name, &user.Email)
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
func GetFollowers(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.follower_id WHERE s.followee_id = $1"
	rows, err := pg.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var followers []User

	for rows.Next() {
		var follower User
		err := rows.Scan(&follower.UserID, &follower.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		followers = append(followers, follower)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(followers)
}

func GetFollowing(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("user_id")
	if userID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	query := "SELECT u.id, u.username FROM users u INNER JOIN subscriptions s ON u.id = s.followee_id WHERE s.follower_id = $1"
	rows, err := pg.DB.Query(query, userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var following []User

	for rows.Next() {
		var followee User
		err := rows.Scan(&followee.UserID, &followee.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		following = append(following, followee)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(following)
}

func SearchUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing search query", http.StatusBadRequest)
		return
	}

	searchQuery := "%" + query + "%"
	query = "SELECT id, name, username FROM users WHERE name ILIKE $1 OR username ILIKE $1"

	rows, err := pg.DB.Query(query, searchQuery)
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

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	userCount, err := services.GetUserCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tweetCount, err := services.GetTweetCount()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	statistics := struct {
		UserCount  int `json:"user_count"`
		TweetCount int `json:"tweet_count"`
	}{
		UserCount:  userCount,
		TweetCount: tweetCount,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statistics)
}
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

func GetUserProfile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	query := "SELECT id, name, bio FROM users_tweeter WHERE id = $1"
	var user User
	err := pg.DB.QueryRow(query, userID).Scan(&user.ID, &user.Name, &user.Bio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Birthday string `json:"birthday"`
		NickName string `json:"nickName"`
		Bio      string `json:"bio"`
		Location string `json:"location"`
	}{
		Name:     user.Name,
		Email:    user.Email,
		Birthday: user.BirthDate,
		NickName: user.Nickname,
		Bio:      user.Bio,
		Location: user.Location,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
func DeleteUserSession(token string) error {
	query := "DELETE FROM user_session WHERE login_token = $1"
	_, err := pg.DB.Exec(query, token)
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
