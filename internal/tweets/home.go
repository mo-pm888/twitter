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
	followingQuery := "SELECT following FROM follower WHERE follower = $1"
	rows, err := pg.DB.Query(followingQuery, userID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()
	var followingIDs []int
	for rows.Next() {
		var following int
		if err = rows.Scan(&following); err != nil {
			log.Println("Error scanning following", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		followingIDs = append(followingIDs, following)
	}

	tweetsQuery := `SELECT tweet_id, text, user_id, created_at,parent_tweet_id,public,only_followers,only_mutual_followers,only_me,retweet FROM tweets WHERE user_id = ANY($1) ORDER BY created_at DESC LIMIT 10`
	followingArray := pq.Array(followingIDs)
	rows, err = pg.DB.Query(tweetsQuery, followingArray)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
	}
	defer rows.Close()

	var tweets []Tweet
	for rows.Next() {
		var tweet Tweet
		if err = rows.Scan(&tweet.TweetID, &tweet.Text, &tweet.Author, &tweet.CreatedAt, &tweet.ParentTweetId, &tweet.Public, &tweet.OnlyFollowers, &tweet.OnlyMutualFollowers, &tweet.OnlyMe, &tweet.Retweet); err != nil {
			log.Println("Error scanning tweets:", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		}
		tweets = append(tweets, tweet)

	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(tweets)
}
