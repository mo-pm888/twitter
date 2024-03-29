package services

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	//Serviceuser "Twitter_like_application/internal/users"
)

type ErrResponse struct {
	Errtext string `json:"errtext"`
}

var (
	emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	emailLen   = 320
)

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

//func UserExists(userID string) bool {
//	query := "SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)"
//	var exists bool
//	err := pg.DB.QueryRow(query, userID).Scan(&exists)
//	if err != nil {
//		return false
//	}
//	return exists
//}

//func IsUserFollowing(currentUserID, targetUserID int) bool {
//	query := "SELECT EXISTS (SELECT 1 FROM subscriptions WHERE user_id = $1 AND target_user_id = $2)"
//	var exists bool
//	err := pg.DB.QueryRow(query, currentUserID, targetUserID).Scan(&exists)
//	if err != nil {
//		return false
//	}
//	return exists
//}

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

//	func GetSubscribedUserIDs(userID string) ([]int, error) {
//		query := "SELECT subscribed_user_id FROM subscriptions WHERE user_id = $1"
//		rows, err := pg.DB.Query(query, userID)
//		if err != nil {
//			return nil, err
//		}
//		defer rows.Close()
//
//		var subscribedUserIDs []int
//
//		for rows.Next() {
//			var subscribedUserID int
//			err := rows.Scan(&subscribedUserID)
//			if err != nil {
//				return nil, err
//			}
//			subscribedUserIDs = append(subscribedUserIDs, subscribedUserID)
//		}
//
//		if err := rows.Err(); err != nil {
//			return nil, err
//		}
//
//		return subscribedUserIDs, nil
//	}
//
//	func GetUserCount() (int, error) {
//		query := "SELECT COUNT(*) FROM users"
//		var count int
//		err := pg.DB.QueryRow(query).Scan(&count)
//		if err != nil {
//			return 0, err
//		}
//		return count, nil
//	}
//
//	func GetTweetCount() (int, error) {
//		query := "SELECT COUNT(*) FROM tweets"
//		var count int
//		err := pg.DB.QueryRow(query).Scan(&count)
//		if err != nil {
//			return 0, err
//		}
//		return count, nil
//	}
func ReturnErr(w http.ResponseWriter, err string, code int) {
	var errj ErrResponse
	errj.Errtext = err
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errj)
}
func CheckEmail(w http.ResponseWriter, email string) {
	if !emailRegex.MatchString(email) {
		ReturnErr(w, "Invalid email format", http.StatusBadRequest)
		return
	}
	if len(email) > 320 {
		ReturnErr(w, "Name exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func CheckPassword(w http.ResponseWriter, password string) {
	passwordRegex := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !passwordRegex.MatchString(password) {
		ReturnErr(w, "Invalid password format", http.StatusBadRequest)
		return
	}
	if len(password) > 100 {
		ReturnErr(w, "Password exceeds maximum length", http.StatusBadRequest)
		return
	}
}
func ReturnJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
}
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLength := big.NewInt(int64(len(charset)))

	randomBytes := make([]byte, length)
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		randomBytes[i] = charset[randomIndex.Int64()]
	}

	return string(randomBytes), nil
}
