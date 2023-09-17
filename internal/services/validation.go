package services

import (
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

var (
	commonWords = []string{"password", "12345678", "87654321", "qwerty123"}
	sequences   = []string{"123", "abc", "xyz"}
	NameRegex   = regexp.MustCompile("^[\\p{L}\\s]+$")
)

func HasDigit(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsDigit(char) {
			return true
		}
	}
	return false
}

func ContainsCommonWord(fl validator.FieldLevel) bool {
	for _, word := range commonWords {
		if strings.Contains(fl.Field().String(), word) {
			return false
		}
	}
	return true
}

func HasNoSequence(fl validator.FieldLevel) bool {
	for _, sequence := range sequences {
		if strings.Contains(fl.Field().String(), sequence) {
			return false
		}
	}
	return true
}

func HasUpper(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if unicode.IsUpper(char) {
			return true
		}
	}
	return false
}

func HasSpecialChar(fl validator.FieldLevel) bool {
	for _, char := range fl.Field().String() {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func CheckDate(fl validator.FieldLevel) bool {
	_, err := time.Parse("2006-01-02", fl.Field().String())
	return err == nil
}
func DateNotAfter(fl validator.FieldLevel) bool {
	date, err := time.Parse("2006-01-02", fl.Field().String())
	if err != nil {
		return false
	}
	currentDate := time.Now()
	return !date.After(currentDate)
}

func CheckNickName(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) < maxNameLength
}
func CheckBio(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) < maxlengthBio
}
func CheckLocation(fl validator.FieldLevel) bool {
	return len(fl.Field().String()) < maxNameLength
}
