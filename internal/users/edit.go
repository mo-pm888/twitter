package users

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type editUserRequest struct {
	createUserRequest
}

func (s *Service) EditProfile(w http.ResponseWriter, r *http.Request) {
	req := editUserRequest{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value("userID").(int)
	err = updateProfile(&req, userID, s.DB)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Profile updated successfully"
	services.ReturnJSON(w, http.StatusOK, message)
}

func updateProfile(req *editUserRequest, userID int, s *sql.DB) error {
	var (
		keys   = []string{}
		values = []any{}
	)
	if err := req.validate(); err != nil {
		return err
	}
	if req.Name != "" {
		values = append(values, req.Name)
		keys = append(keys, " name = $"+strconv.Itoa(len(keys)+1))
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return err
		}
		values = append(values, string(hashedPassword))
		keys = append(keys, " password = $"+strconv.Itoa(len(keys)+1))
	}
	if req.Email != "" {
		values = append(values, req.Email)
		keys = append(keys, " email = $"+strconv.Itoa(len(keys)+1))
	}
	if req.BirthDate != "" {
		values = append(values, req.BirthDate)
		keys = append(keys, " birthdate = $"+strconv.Itoa(len(keys)+1))
	}

	if req.Nickname != "" {
		values = append(values, req.Nickname)
		keys = append(keys, " nickname = $"+strconv.Itoa(len(keys)+1))
	}
	if req.Bio != "" {
		values = append(values, req.Bio)
		keys = append(keys, " bio = $"+strconv.Itoa(len(keys)+1))
	}

	if req.Location != "" {
		values = append(values, req.Location)
		keys = append(keys, " location = $"+strconv.Itoa(len(keys)+1))

	}
	values = append(values, userID)
	keyString := strings.Join(keys, ", ")
	query := fmt.Sprintf("UPDATE users_tweeter SET %s WHERE id = $%d", keyString, len(values))
	_, err := s.Exec(query, values...)
	if err != nil {
		return err
	}
	return err
}

func (s editUserRequest) validate() error {
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
	if err := v.RegisterValidation("hasCommonWord", services.HasCommonWord); err != nil {
		return err
	}
	if err := v.RegisterValidation("hasDigit", services.HasDigit); err != nil {
		return err
	}
	if err := v.RegisterValidation("date", services.CheckDate); err != nil {
		return err
	}
	if err := v.RegisterValidation("dateAfter", services.CheckDateAfter); err != nil {
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
