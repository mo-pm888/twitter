package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createUserRequest_validate(t *testing.T) {
	r := createUserRequest{
		Name:     "!!!",
		Email:    "asdasd",
		Password: "asdasdfasdfsadfas",
	}

	err := r.validate()

	assert.Error(t, err)
	t.Log(err.Error())
}
