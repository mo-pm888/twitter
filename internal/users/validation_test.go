package users

import (
	"fmt"
	"testing"
	"time"

	"Twitter_like_application/internal/services"

	"github.com/go-playground/validator/v10"
)

func TestCheckPassword(t *testing.T) {
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		Password string `validate:"omitempty,checkPassword"`
		Expected bool
	}{
		{"ValidPassword819!", true},
		{"Sh1!", false},
		{"invalidnouppercase1!", false},
		{"InvalidNoSpecialChar1", false},
		{"InvalidNoSpecialChar1123213wcewewc!!!BNFKJDBFLKJWBFL:WEBFoij3oihjoi3hr3h9848fhdsubfkfjsdbgkjsbnflsdbnfu3b289fbl;sdbf3", false},
		{"ValidPassword123!", false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}

func TestCheckName(t *testing.T) {
	randomName, err := services.GenerateRandomString(200)
	if err != nil {
		fmt.Println(err)
		return
	}
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		Name     string `validate:"omitempty,checkName"`
		Expected bool
	}{
		{"Alex", true},
		{randomName, false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}

func TestCheckDate(t *testing.T) {
	currentDate := time.Now()
	tomorrow := currentDate.AddDate(0, 0, 1)
	tomorrowFormatted := tomorrow.Format("2006-01-02")
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		BirthDate string `validate:"omitempty,checkDate"`
		Expected  bool
	}{
		{"1988-08-08", true},
		{"442343-222-3", false},
		{tomorrowFormatted, false},
		{"dfgsdfd", false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}
func TestCheckNickName(t *testing.T) {
	randomNick, err := services.GenerateRandomString(200)
	if err != nil {
		fmt.Println(err)
		return
	}
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		NickName string `validate:"omitempty,checkNickname"`
		Expected bool
	}{
		{"bunin", true},
		{randomNick, false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}
func TestCheckBio(t *testing.T) {
	randomBio, err := services.GenerateRandomString(500)
	if err != nil {
		fmt.Println(err)
		return
	}
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		Bio      string `validate:"omitempty,checkBio"`
		Expected bool
	}{
		{"test bio right", true},
		{randomBio, false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}
func TestCheckEmail(t *testing.T) {
	testValid := &UserValid{
		validate: validator.New(),
		validErr: make(map[string]string),
	}
	testCases := []struct {
		Email    string `validate:"omitempty,email"`
		Expected bool
	}{
		{"test@mail.com", true},
		{"dfffff", false},
	}
	if err := RegisterUsersValidations(testValid); err != nil {
		fmt.Println(err)
		return
	}
	for i := 0; i < len(testCases); i++ {
		err := testValid.validate.Struct(testCases[i])
		if (err == nil) == testCases[i].Expected {
			t.Logf("Test case %d passed", i+1)
		} else {
			t.Errorf("Test case %d failed, expected %v, got error: %v", i+1, testCases[i].Expected, err)
		}
	}
}
