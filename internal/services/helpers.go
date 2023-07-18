package services

import (
	Postgresql "Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/tweets"
	"time"

	//Serviceuser "Twitter_like_application/internal/users"

	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ErrResponse struct {
	Errtext string `json:"errtext"`
}

func GenerateResetToken() string {
	const resetTokenLength = 32
	tokenBytes := make([]byte, resetTokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(tokenBytes)
}

func ConvertStringToNumber(str string) (int, error) {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0, err
	}
	return num, nil
}

func UserExists(userID string) bool {
	query := "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)"
	var exists bool
	err := Postgresql.DB.QueryRow(query, userID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func IsUserFollowing(currentUserID, targetUserID int) bool {
	query := "SELECT EXISTS (SELECT 1 FROM subscriptions WHERE user_id = $1 AND target_user_id = $2)"
	var exists bool
	err := Postgresql.DB.QueryRow(query, currentUserID, targetUserID).Scan(&exists)
	if err != nil {
		return false
	}
	return exists
}

func GetCurrentUserID(r *http.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return "", errors.New("No session cookie found")
	} else if err != nil {
		return "", err
	}

	userID, err := ExtractUserIDFromSessionCookie(cookie.Value)

	return userID, nil
}
func ExtractUserIDFromSessionCookie(cookieValue string) (string, error) {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader("GET / HTTP/1.0\r\nCookie: session=" + cookieValue + "\r\n\r\n")))
	if err != nil {
		return "", err
	}

	cookie, err := req.Cookie("session")
	if err != nil {
		return "", err
	}

	userID := cookie.Value
	return userID, nil
}
func GetSubscribedUserIDs(userID string) ([]int, error) {
	query := "SELECT subscribed_user_id FROM subscriptions WHERE user_id = $1"
	rows, err := Postgresql.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribedUserIDs []int

	for rows.Next() {
		var subscribedUserID int
		err := rows.Scan(&subscribedUserID)
		if err != nil {
			return nil, err
		}
		subscribedUserIDs = append(subscribedUserIDs, subscribedUserID)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscribedUserIDs, nil
}
func GetUserCount() (int, error) {
	query := "SELECT COUNT(*) FROM users"
	var count int
	err := Postgresql.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func GetTweetCount() (int, error) {
	query := "SELECT COUNT(*) FROM tweets"
	var count int
	err := Postgresql.DB.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
func ReturnErr(w http.ResponseWriter, err string, code int) {
	var errj ErrResponse
	errj.Errtext = err
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errj)
}
func CreatNewTweet(newTweet tweets.CreatNewTweet, r *http.Request, userID int, tweetTo string) error {
	query := `INSERT INTO tweets (user_id, text, created_at, public, only_followers, only_mutual_followers, only_me,reply_to)
		VALUES ($1, $2, $3, $4, $5, $6, $7,&8) RETURNING tweet_id`
	err := Postgresql.DB.QueryRowContext(r.Context(), query, userID, newTweet.Text, time.Now(), newTweet.Public, newTweet.OnlyFollowers, newTweet.OnlyMutualFollowers, newTweet.OnlyMe, tweetTo).Scan(&newTweet.TweetID)
	if err != nil {
		return err
	}
	return nil
}
func ReturnJSON(w http.ResponseWriter, statusCode int, message string) {
	response := map[string]interface{}{
		"status":  "success",
		"message": message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
}
