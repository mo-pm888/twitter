package tweets

import (
	"encoding/json"
	"net/http"
	"time"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

const maxLengthTweet = 400

type CreatNewTweet struct {
	TweetID             int
	Text                string `json:"text" validate:"required,checkTweetText"`
	CreatedAt           time.Time
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}

type createTweetRequest struct {
	ID                  int    `json:"id"`
	Text                string `json:"text" validate:"required,text"`
	CreatedAt           time.Time
	Public              bool `json:"public"`
	OnlyFollowers       bool `json:"only_followers"`
	OnlyMutualFollowers bool `json:"only_mutual_followers"`
	OnlyMe              bool `json:"only_me"`
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req := createTweetRequest{}
	userID := r.Context().Value("userID").(int)
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := req.validate(); err != nil {
		services.ReturnErr(w, err, http.StatusBadRequest)
		return
	}

	query := `INSERT INTO tweets (user_id, text, created_at, public, only_followers, only_mutual_followers, only_me)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING tweet_id`
	err = s.DB.QueryRowContext(r.Context(), query, userID, req.Text, time.Now(), req.Public, req.OnlyFollowers, req.OnlyMutualFollowers, req.OnlyMe).Scan(&req.ID)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return

	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(req)

	return
}
func (s createTweetRequest) validateText(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) > maxLengthTweet
}

func (s createTweetRequest) validate() error {
	v := validator.New()
	if err := v.RegisterValidation("text", s.validateText); err != nil {
		return err
	}
	return v.Struct(s)
}
