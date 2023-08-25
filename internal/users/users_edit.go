package users

import (
	"Twitter_like_application/internal/database/pg"
	"Twitter_like_application/internal/services"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
)

type EditUserRequest struct {
	ID        int    `json:"-"`
	Name      string `json:"name" validate:"omitempty,checkName"`
	Password  string `json:"password" validate:"omitempty,checkPassword"`
	Email     string `json:"email" validate:"omitempty,email"`
	BirthDate string `json:"birthdate" validate:"omitempty,checkDate"`
	Nickname  string `json:"nickname" validate:"omitempty,checkNickname"`
	Bio       string `json:"bio" validate:"omitempty,checkBio"`
	Location  string `json:"location" validate:"omitempty,checkLocation"`
}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	userValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	updatedProfile := EditUserRequest{}
	if err := RegisterUsersValidations(userValid); err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := r.Context().Value("userID").(int)
	err = updateProfile(&updatedProfile, userID, userValid)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return
	}

	message := "Profile updated successfully"
	services.ReturnJSON(w, http.StatusOK, message)
}

func updateProfile(updatedProfile *EditUserRequest, userID int, v *UserValid) error {
	var (
		hashedPassword []byte
		keys           = []string{}
		values         = []any{}
	)
	err := v.validate.Struct(updatedProfile)
	if err != nil {
		return err
	}
	if updatedProfile.Name != "" {
		values = append(values, updatedProfile.Name)
		keys = append(keys, " name = $"+strconv.Itoa(len(keys)+1))
	}
	if updatedProfile.Password != "" {
		hashedPassword, err = bcrypt.GenerateFromPassword([]byte(updatedProfile.Password), bcrypt.DefaultCost)
		if err != nil {
			fmt.Println(err)
			return err
		}
		values = append(values, string(hashedPassword))
		keys = append(keys, " password = $"+strconv.Itoa(len(keys)+1))
	}
	if updatedProfile.Email != "" {
		values = append(values, updatedProfile.Email)
		keys = append(keys, " email = $"+strconv.Itoa(len(keys)+1))
	}
	if updatedProfile.BirthDate != "" {
		values = append(values, updatedProfile.BirthDate)
		keys = append(keys, " birthdate = $"+strconv.Itoa(len(keys)+1))
	}

	if updatedProfile.Nickname != "" {
		values = append(values, updatedProfile.Nickname)
		keys = append(keys, " nickname = $"+strconv.Itoa(len(keys)+1))
	}
	if updatedProfile.Bio != "" {
		values = append(values, updatedProfile.Bio)
		keys = append(keys, " bio = $"+strconv.Itoa(len(keys)+1))
	}

	if updatedProfile.Location != "" {
		values = append(values, updatedProfile.Location)
		keys = append(keys, " location = $"+strconv.Itoa(len(keys)+1))

	}
	values = append(values, userID)
	keyString := strings.Join(keys, ", ")
	query := fmt.Sprintf("UPDATE users_tweeter SET %s WHERE id = $%d", keyString, len(values))
	_, err = pg.DB.Exec(query, values...)
	if err != nil {
		return err
	}
	return err
}
