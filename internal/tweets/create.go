package tweets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

type createTweetRequest struct {
	Text string `json:"text" validate:"required,max=400"`
	Visibility
}

var DefaultVisibility = Visibility{
	Public: true,
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req := createTweetRequest{}
	userID := r.Context().Value("u`serID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = s.CreateNewTweet(req, r.Context(), userID); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

	return
}
func (s *Service) CreateNewTweet(req createTweetRequest, ctx context.Context, parentID int) error {
	if err := req.validate(); err != nil {
		return err
	}
	userID := ctx.Value("userID").(int)
	if req.defaultVisibilities() {
		req.Visibility = DefaultVisibility
	}
	query := `INSERT INTO tweets (user_id, text, created_at,parent_tweet_id, public, only_followers, only_mutual_followers, only_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7,$8)`
	if err := s.DB.QueryRowContext(ctx, query, userID, req.Text, time.Now(), parentID, req.Public, req.OnlyFollowers, req.OnlyMutualFollowers, req.OnlyMe).Err(); err != nil {
		return err
	}
	return nil

}

func (s createTweetRequest) validate() error {
	v := validator.New()
	if !s.Visibility.isValid() {
		return errors.New("visibility error, many visual arguments")
	}
	return v.Struct(s)
}
