package tweets

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"strings"
)

var (
	maxLenghtTweet = 400
)

func (v *TweetValid) Error() string {
	var pairs []string
	for k, v := range v.ValidErr {
		pairs = append(pairs, fmt.Sprintf("%s: %s", k, v))
	}

	result := strings.Join(pairs, "; ")
	return result
}

func CheckTweetText(fl validator.FieldLevel, v *TweetValid) bool {
	text := fl.Field().String()
	if len(text) > maxLenghtTweet {
		v.ValidErr["name"] += "long text,"
	}
	return true
}
func RegisterTweetValidations(tweetValid *TweetValid) error {
	err := tweetValid.Validate.RegisterValidation("checkTweetText", func(fl validator.FieldLevel) bool {
		return CheckTweetText(fl, tweetValid)
	})
	if err != nil {
		return err
	}
	return nil
}
