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

	userID := r.Context().Value("userID").(string)
	userIDint, err := strconv.Atoi(userID)
	offset := (page - 1) * perPage
	limit := perPage

	followingQuery := "SELECT following FROM follower WHERE follower = $1"
	rows, err := pg.DB.Query(followingQuery, userIDint)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {

		}
	}(rows)

	var followingIDs []int
	for rows.Next() {
		var following int
		if err = rows.Scan(&following); err != nil {
			log.Println("Error scanning following:", err)
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		followingIDs = append(followingIDs, following)
	}

	tweetsQuery := `
	SELECT tweet_id, text, user_id, created_at,parent_tweet_id,public,
	only_followers,only_mutual_followers,only_me,retweet
	FROM tweets
	WHERE user_id = ANY($1)
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
`
	followingArray := pq.Array(followingIDs)
	rows, err = pg.DB.Query(tweetsQuery, followingArray, limit, offset)
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
