package services

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const (
	maxNameLength  = 100
	minNameLength  = 8
	maxlengthBio   = 400
	maxLenghtTweet = 400
)

type Services struct {
	Validate *validator.Validate
	ValidErr map[string]string
}
type NameVal struct {
	short    bool
	long     bool
	realName bool
}

func New() (*Services, error) {
	userValid := &Services{
		Validate: validator.New(),
		ValidErr: make(map[string]string),
	}
	if err := userValid.RegisterValidations(userValid); err != nil {
		return nil, err
	}
	return userValid, nil
}

func (s *Services) Error() string {
	var pairs []string
	for k, v := range s.ValidErr {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
	}

	return strings.Join(pairs, "; ")
}

func (s *Services) RegisterValidations(userValid *Services) error {
	err := userValid.Validate.RegisterValidation("checkDate", func(fl validator.FieldLevel) bool {
		return s.CheckDate(fl)
	})
	if err != nil {
		return err
	}
	err = userValid.Validate.RegisterValidation("checkNickname", func(fl validator.FieldLevel) bool {
		return s.CheckNickName(fl)
	})
	if err != nil {
		return err
	}
	err = userValid.Validate.RegisterValidation("checkBio", func(fl validator.FieldLevel) bool {
		return s.CheckBio(fl)
	})
	if err != nil {
		return err
	}
	err = userValid.Validate.RegisterValidation("checkLocation", func(fl validator.FieldLevel) bool {
		return s.CheckLocation(fl)
	})
	if err != nil {
		return err
	}
	err = userValid.Validate.RegisterValidation("email", func(fl validator.FieldLevel) bool {
		return s.CheckEmail(fl)
	})
	if err != nil {
		return err
	}
	err = userValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return s.CheckTweetText(fl)
	})
	if err != nil {
		return err
	}
	return nil
}

var (
	commonWords = []string{"password", "12345678", "87654321", "qwerty123"}
	sequences   = []string{"123", "abc", "xyz"}
	NameRegex   = regexp.MustCompile("^[\\p{L}\\s]+$")
)

func HasDigit(password string) bool {
	for _, char := range password {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func HasCommonWord(password string) bool {
	for _, word := range commonWords {
		if strings.Contains(password, word) {
			return true
		}
	}
	return false
}

func HasSequence(password string) bool {
	for _, sequence := range sequences {
		if strings.Contains(password, sequence) {
			return true
		}
	}
	return false
}
func HasUpper(password string) bool {
	for _, char := range password {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}
func HasSpecialChar(password string) bool {
	for _, char := range password {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}
func CheckPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	errorsMap := make(map[string]string)
	if !HasUpper(password) {
		errorsMap["password"] += "uppercase,"
	}
	if !HasSpecialChar(password) {
		errorsMap["password"] += "special character,"
	}
	if !HasDigit(password) {
		errorsMap["password"] += "digit,"
	}
	if HasSequence(password) {
		errorsMap["password"] += "sequence,"
	}
	if HasCommonWord(password) {
		errorsMap["password"] += "common word,"
	}
	return len(errorsMap) == 0
}

func (s *Services) CheckDate(fl validator.FieldLevel) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		s.ValidErr["data"] += "incorrect date,"
		return false
	}
	currentDate := time.Now()
	if date.After(currentDate) {
		s.ValidErr["data"] += "date is after current date"
		return true
	}

	return true
}

func (s *Services) CheckNickName(fl validator.FieldLevel) bool {
	nickname := fl.Field().String()
	if len(nickname) > maxNameLength {
		s.ValidErr["nickname"] = "long"
	}
	return true

}
func (s *Services) CheckBio(fl validator.FieldLevel) bool {
	if len(fl.Field().String()) > maxlengthBio {
		s.ValidErr["bio"] = "long"
	}
	return true

}
func (s *Services) CheckLocation(fl validator.FieldLevel) bool {
	location := fl.Field().String()
	if len(location) > maxNameLength {
		s.ValidErr["location"] = "long"
	}
	return true

}
func (s *Services) CheckEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	_, err := mail.ParseAddress(email)
	if err != nil {
		s.ValidErr["email"] = "not correct email"
		return false
	}
	return true
}
func (s *Services) CheckTweetText(fl validator.FieldLevel) bool {
	text := fl.Field().String()
	if len(text) > maxLenghtTweet {
		s.ValidErr["name"] += "long text,"
	}
	return true
}
