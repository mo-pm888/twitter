package users

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"

	"Twitter_like_application/internal/services"
	Tweets "Twitter_like_application/internal/tweets"

	"github.com/go-playground/validator/v10"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                 int
	Name               string `json:"name"`
	Password           string `json:"password"`
	Email              string `json:"email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birthdate"`
	Nickname           string `json:"nickname"`
	Bio                string `json:"bio"`
	Location           string `json:"location"`
	Tweets.Tweet
}

type createUserRequest struct {
	Name      string `json:"name" validate:"required,max=100,checkName"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=8,max=100,hasUpper,hasSpecialChar,hasSequence,hasCommonWord,hasDigit"`
	BirthDate string `json:"birthdate" validate:"required,date,dateAfter"`
	Nickname  string `json:"nickname" validate:"omitempty,nickName"`
	Bio       string `json:"bio" validate:"omitempty,bio"`
	Location  string `json:"location" validate:"omitempty,location"`
}

func (s *Service) Create(w http.ResponseWriter, r *http.Request) {
	req := &createUserRequest{}
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = req.validate(); err != nil {
		services.ReturnErr(w, err, http.StatusBadRequest)
		return
	}
	query := `SELECT id FROM users_tweeter WHERE email = $1`
	var existingUserID int
	err = s.DB.QueryRow(query, req.Email).Scan(&existingUserID)
	if err == nil {
		services.ReturnErr(w, "The user has already existed with this email ", http.StatusBadRequest)
		return
	} else if err != sql.ErrNoRows {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	req.Password = string(hashedPassword)
	query = `INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	err = s.DB.QueryRow(query, req.Name, req.Password, req.Email, req.Nickname, req.Location, req.Bio, req.BirthDate).Err()
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			services.ReturnErr(w, "This user is already added", http.StatusBadRequest)
			return
		}
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	EmailVerificationToken(req.Email)

	services.ReturnJSON(w, http.StatusCreated, "new user was created")
}

func (s createUserRequest) validateName(fl validator.FieldLevel) bool {
	return services.NameRegex.MatchString(fl.Field().String())
}

func (s createUserRequest) validateEmail(fl validator.FieldLevel) bool {
	_, err := mail.ParseAddress(fl.Field().String())
	return err == nil
}

func (s createUserRequest) validate() error {
	v := validator.New()
	if err := v.RegisterValidation("checkName", s.validateName); err != nil {
		return err
	}
	if err := v.RegisterValidation("email", s.validateEmail); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasUpper", services.HasUpper); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasSpecialChar", services.HasSpecialChar); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasSequence", services.HasNoSequence); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasCommonWord", services.ContainsCommonWord); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasDigit", services.HasDigit); err != nil {
		return err
	}
	if err := v.RegisterValidation("date", services.CheckDate); err != nil {
		return err
	}
	if err := v.RegisterValidation("dateAfter", services.DateNotAfter); err != nil {
		return err
	}
	if err := v.RegisterValidation("nickName", services.CheckNickName); err != nil {
		return err
	}
	if err := v.RegisterValidation("bio", services.CheckBio); err != nil {
		return err
	}
	if err := v.RegisterValidation("location", services.CheckLocation); err != nil {
		return err
	}
	return v.Struct(s)
}
