package users

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"Twitter_like_application/internal/services"
	Tweets "Twitter_like_application/internal/tweets"

	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 int
	Name               string `json:"name" validate:"omitempty"`
	Password           string `json:"password" validate:"omitempty"`
	Email              string `json:"email" validate:"omitempty,email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birthdate" validate:"omitempty"`
	Nickname           string `json:"nickname" validate:"omitempty"`
	Bio                string `json:"bio" validate:"omitempty"`
	Location           string `json:"location" validate:"omitempty"`
	Tweets.Tweet
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	newUser := &User{}
	err := json.NewDecoder(r.Body).Decode(newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	query := `SELECT id FROM users_tweeter WHERE email = $1`
	var existingUserID int
	err = s.DB.QueryRow(query, newUser.Email).Scan(&existingUserID)
	if err == nil {
		services.ReturnErr(w, "The user has already existed with this email ", http.StatusBadRequest)
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
	err = s.DB.QueryRow(query, newUser.Name, newUser.Password, newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).Scan(&newUser.ID)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			services.ReturnErr(w, "This user is already added", http.StatusBadRequest)
			return
		}
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newUser.EmailToken = CheckEmail(newUser)

	services.ReturnJSON(w, http.StatusCreated, "new user was created")
}
