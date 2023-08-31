package tweets

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

func (v *TweetValid) Error() string {
	var pairs []string
	for k, v := range v.ValidErr {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
	}

	result := strings.Join(pairs, "; ")
	return result
}

func (s *Service) CheckTweetText(fl validator.FieldLevel, v *TweetValid) bool {
	text := fl.Field().String()
	if len(text) > s.TweetLength {
		v.ValidErr["tweet"] += "long text,"
	}
	return true
}
func (s *Service) RegisterTweetValidations(tweetValid *TweetValid) error {
	err := tweetValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return s.CheckTweetText(fl, tweetValid)
	})
	if err != nil {
		return err
	}
	return nil
}
