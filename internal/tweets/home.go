package tweets

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"github.com/lib/pq"
	"log"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	followsQuery := "SELECT follower_id FROM followers_subscriptions WHERE subscription_id = $1"
	rows, err := pg.DB.Query(followsQuery, userID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()
	var followers []int
	for rows.Next() {
		var followee int
		if err = rows.Scan(&followee); err != nil {
			log.Println("Error scanning followers:", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		followers = append(followers, followee)
	}

	tweetsQuery := `SELECT tweet_id, text, user_id, created_at FROM tweets WHERE user_id = ANY($1) ORDER BY created_at DESC`
	followersArray := pq.Array(followers)
	rows, err = pg.DB.Query(tweetsQuery, followersArray)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var tweets []Tweet
	for rows.Next() {
		var tweet Tweet
		if err = rows.Scan(&tweet.TweetID, &tweet.Text, &tweet.Author, &tweet.CreatedAt); err != nil {
			log.Println("Error scanning tweets:", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		tweets = append(tweets, tweet)

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweets)
}
