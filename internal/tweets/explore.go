package tweets

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
)

func ExploreRandom(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		services.ReturnErr(w, "Invalid page parameter", http.StatusBadRequest)
		return
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil {
		services.ReturnErr(w, "Invalid per_page parameter", http.StatusBadRequest)
		return
	}

	offset := (page - 1) * perPage
	limit := perPage

	tweetsQuery := `
		SELECT tweet_id, text, user_id, created_at, parent_tweet_id, public,
		only_followers, only_mutual_followers, only_me, retweet
		FROM tweets
		ORDER BY RANDOM()
		LIMIT $1 OFFSET $2
	`

	rows, err := pg.DB.Query(tweetsQuery, limit, offset)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var tweets []Tweet
	for rows.Next() {
		var tweet Tweet
		if err = rows.Scan(&tweet.TweetID, &tweet.Text, &tweet.UserID, &tweet.CreatedAt,
			&tweet.ParentTweetId, &tweet.Public, &tweet.OnlyFollowers,
			&tweet.OnlyMutualFollowers, &tweet.OnlyMe, &tweet.Retweet); err != nil {
			log.Println("Error scanning tweets:", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tweets = append(tweets, tweet)
	}

	w.WriteHeader(http.StatusOK)
	services.ReturnJSON(w, http.StatusOK, tweets)
}
