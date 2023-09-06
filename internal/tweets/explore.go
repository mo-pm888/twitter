package tweets

import (
	"net/http"
	"strconv"

	"Twitter_like_application/internal/services"
)

const (
	badPage    = "invalid page parameter"
	defPage    = "1"
	defPerPage = "1"
)

func (s *Service) Explore(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	perPageStr := r.URL.Query().Get("per_page")
	if pageStr == "" && perPageStr == "" {
		pageStr = defPage
		perPageStr = defPerPage
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 || page > 100 {
		services.ReturnErr(w, badPage, http.StatusBadRequest)
		return
	}

	perPage, err := strconv.Atoi(perPageStr)
	if err != nil || perPage <= 0 || perPage > 100 {
		services.ReturnErr(w, badPage, http.StatusBadRequest)
		return
	}

	offset := (page - 1) * perPage
	limit := perPage

	userID, ok := r.Context().Value("userID").(int)
	if ok && userID != 0 {
		query := `
SELECT tweet_id, text, user_id, created_at, parent_tweet_id, public,
       only_followers, only_mutual_followers, only_me
FROM tweets
WHERE user_id != $1 
  AND (
        user_id NOT IN (
            SELECT following FROM follower WHERE follower = $1
        ) 
        OR retweet != 0
      )
  AND public = true
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

`
		rows, err := s.DB.Query(query, userID, limit, offset)
		if err != nil {
			services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var tweets []Tweet
		for rows.Next() {
			var tweet Tweet
			if err := rows.Scan(&tweet.TweetID, &tweet.Text, &tweet.UserID, &tweet.CreatedAt,
				&tweet.ParentTweetId, &tweet.Public, &tweet.OnlyFollowers,
				&tweet.OnlyMutualFollowers, &tweet.OnlyMe, &tweet.Retweet); err != nil {
				services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
				return
			}
			tweets = append(tweets, tweet)
		}

		w.WriteHeader(http.StatusOK)
		services.ReturnJSON(w, http.StatusOK, tweets)
	}
}
