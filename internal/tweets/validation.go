package tweets

import (
	"fmt"
	"strconv"
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

func CheckTweetText(fl validator.FieldLevel, v *TweetValid, s *Service) bool {
	maxLengthTweet, err := strconv.Atoi(s.TweetLength)
	if err != nil {
		return false
	}
	text := fl.Field().String()
	fmt.Println(maxLengthTweet)
	if len(text) > maxLengthTweet {
		v.ValidErr["tweet"] += "long text,"
	}
	return true
}
func RegisterTweetValidations(tweetValid *TweetValid, s *Service) error {
	err := tweetValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return CheckTweetText(fl, tweetValid, s)
	})
	if err != nil {
		return err
	}
	return nil
}
