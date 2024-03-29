package users

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
	maxNameLength = 100
	minNameLength = 8
	maxlengthBio  = 400
)

var (
	commonWords = []string{"password", "12345678", "87654321", "qwerty123"}
	sequences   = []string{"123", "abc", "xyz"}
	nameRegex   = regexp.MustCompile("^[\\p{L}\\s]+$")
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
func CheckPassword(fl validator.FieldLevel, v *UserValid) bool {
	password := fl.Field().String()
	if len(password) < minNameLength {
		v.validErr["password"] += "short,"
	}
	if len(password) > maxNameLength {
		v.validErr["password"] += "long,"
	}

	if !HasUpper(password) {
		v.validErr["password"] += "uppercase,"
	}
	if !HasSpecialChar(password) {
		v.validErr["password"] += "special character,"
	}
	if !HasDigit(password) {
		v.validErr["password"] += "digit,"
	}
	if HasSequence(password) {
		v.validErr["password"] += "sequence,"
	}
	if HasCommonWord(password) {
		v.validErr["password"] += "common word,"
	}
	return len(v.validErr) == 0
}

func CheckDate(fl validator.FieldLevel, v *UserValid) bool {
	dateStr := fl.Field().String()
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		v.validErr["data"] += "incorrect date,"
		return false
	}
	currentDate := time.Now()
	if date.After(currentDate) {
		v.validErr["data"] += "date is after current date"
		return false
	}

	return true
}

func CheckName(fl validator.FieldLevel, v *UserValid) bool {
	name := fl.Field().String()
	u := NameVal{}
	if len(name) > maxNameLength {
		v.validErr["name"] += "long name,"
		u.long = true
	}
	match := nameRegex.MatchString(name)
	if match == false {
		v.validErr["name"] += "digit or special character,"
		u.realName = true
	}

	if u.long || u.realName {
		return false
	}
	return true
}

func CheckNickName(fl validator.FieldLevel, v *UserValid) bool {
	nickname := fl.Field().String()
	if len(nickname) > maxNameLength {
		v.validErr["nickname"] = "long"
		return false
	}
	return true

}
func CheckBio(fl validator.FieldLevel, v *UserValid) bool {
	if len(fl.Field().String()) > maxlengthBio {
		v.validErr["bio"] = "long"
		return false
	}
	return true

}
func CheckLocation(fl validator.FieldLevel, v *UserValid) bool {
	location := fl.Field().String()
	if len(location) > maxNameLength {
		v.validErr["location"] = "long"
	}
	return true

}
func CheckEmailVal(fl validator.FieldLevel, v *UserValid) bool {
	email := fl.Field().String()
	_, err := mail.ParseAddress(email)
	fmt.Println(err)
	if err != nil {
		v.validErr["email"] = "not correct email"
		return false
	}
	return true
}
