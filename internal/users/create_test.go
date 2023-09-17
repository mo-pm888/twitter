package users

import (
	"testing"
)

func Test_createUserRequest_validate(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
		}

		if err := r.validate(); err != nil {
			t.Errorf("expect: err==nil, got: %s", err)
		}
	})
	t.Run("name_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Name' Error:Field validation for 'Name' failed on the 'checkName' tag"
		r := createUserRequest{
			Name:      "111",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
		}
		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("name_length_ok", func(t *testing.T) {
		r := createUserRequest{
			Name:      "nameOK",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
		}
		if err := r.validate(); err != nil {
			t.Errorf("expect: err==nil, got: %s", err)
		}
	})
	t.Run("name_length_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Name' Error:Field validation for 'Name' failed on the 'max' tag"
		r := createUserRequest{
			Name:      "vmrvmwrpvmwpormvwpormvwpomrpohwrponhptowmhpowmthopwmthponwtpohnqpohnqpo5nhpownhpwonthpownthponwtphonwpothnpwonhpownthponwptohnwponthpownmtphonwpthnwponhtwpotnhpwonthpownthpownthpowntphonwpothnwptohn",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
		}
		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("email_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Email' Error:Field validation for 'Email' failed on the 'email' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd3",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("birthDate_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.BirthDate' Error:Field validation for 'BirthDate' failed on the 'date' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07f",
			Password:  "dgfghheeeDF1@",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("password_hasUpper_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Password' Error:Field validation for 'Password' failed on the 'hasUpper' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dfdsfsfsf12!",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("password_hasSpecialChar_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Password' Error:Field validation for 'Password' failed on the 'hasSpecialChar' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dfdsfsfsfF12",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("password_hasDigit_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Password' Error:Field validation for 'Password' failed on the 'hasDigit' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dfdsfsfsfF!",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("password_hasSequence_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Password' Error:Field validation for 'Password' failed on the 'hasSequence' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dfdsfsfsfF!123",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("bio_ok", func(t *testing.T) {
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
			Bio:       "bioOK",
		}

		if err := r.validate(); err != nil {
			t.Errorf("expect: err==nil, got: %s", err)
		}
	})
	t.Run("bio_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Bio' Error:Field validation for 'Bio' failed on the 'max' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
			Bio:       "PGOJAWPEORGJOPERJGPOjergpojwpoergjporjgporjapogjaproegjpoarjgpojrpogjporjgpojrobioOKdakngapnrmopnqmroinvboreinboairnbaoinboirnboainrboianboinaoibnaoirnboianrboinaoibnarnboainrbioanbionarobinaroibnoairnboanrobinaroibnaoirnboianobianeroibnaoe;ribnaoierbnaoiernboaiernboiaenrboinaoribnoiarenboiarnboinarboinaroibnaiornboainboianweIOGHWoeihgoiwhGOIHWOIGHOIWH4GOIhw4oighoIHRGOihgoihOI4GHoigh4oiewgWEGwegwEGwegWEGewgwGwgwGwgWG",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}
	})
	t.Run("location_ok", func(t *testing.T) {
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
			Location:  "locationOk",
		}

		if err := r.validate(); err != nil {
			t.Errorf("expect: err==nil, got: %s", err)
		}
	})
	t.Run("location_fail", func(t *testing.T) {
		expectedErrorSubstring := "Key: 'createUserRequest.Location' Error:Field validation for 'Location' failed on the 'max' tag"
		r := createUserRequest{
			Name:      "kli",
			Email:     "asdasd@mail.ru",
			BirthDate: "1987-12-07",
			Password:  "dgfghheeeDF1@",
			Location:  "PGOJAWPEORGJOPERJGPOjergpojwpoergjporjgporjapogjaproegjpoarjgpojrpogjporjgpojrobioOKdakngapnrmopnqmroinvboreinboairnbaoinboirnboainrboianboinaoibnaoirnboianrboinaoibnarnboainrbioanbionarobinaroibnoairnboanrobinaroibnaoirnboianobianeroibnaoe;ribnaoierbnaoiernboaiernboiaenrboinaoribnoiarenboiarnboinarboinaroibnaiornboainboianweIOGHWoeihgoiwhGOIHWOIGHOIWH4GOIhw4oighoIHRGOihgoihOI4GHoigh4oiewgWEGwegwEGwegWEGewgwGwgwGwgWG",
		}

		err := r.validate()
		if err.Error() != expectedErrorSubstring {
			t.Errorf("Expected error message: %s, Actual error message: %s", expectedErrorSubstring, err.Error())
		}

	})

}
