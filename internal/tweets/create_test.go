package tweets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createTweetRequest_validate(t *testing.T) {
	r := createTweetRequest{
		Text: "gfgfgfgf",
	}

	err := r.validate()

	assert.Error(t, err)
	t.Log(err.Error())
}
