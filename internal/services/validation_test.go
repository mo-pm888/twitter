package services

import (
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
)

type TestInThePastStruct struct {
	Date string `validate:"InThePast"`
}
type TestCommonPasswordStruct struct {
	CommonPassword string `validate:"hasCommonWord"`
}
type TestSequenceStruct struct {
	SequencePassword string `validate:"hasSequence"`
}
type TestDigitStruct struct {
	DigitPassword string `validate:"hasDigit"`
}
type TestSpecialCharStruct struct {
	SpecialCharPassword string `validate:"hasSpecialChar"`
}

func TestInThePast(t *testing.T) {
	v := validator.New()

	if err := v.RegisterValidation("InThePast", InThePast); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		data TestInThePastStruct
		want bool
	}{
		{
			name: "Date in the past",
			data: TestInThePastStruct{Date: "2022-01-01"},
			want: true,
		},
		{
			name: "Date in the future",
			data: TestInThePastStruct{Date: "2024-01-01"},
			want: false,
		},
		{
			name: "Invalid date",
			data: TestInThePastStruct{Date: "invalid-date"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Struct(tt.data); (err == nil) != tt.want {
				t.Errorf("InThePast() = %v, want %v", (err == nil), tt.want)
			}
		})
	}
}

func TestContainsCommonWord(t *testing.T) {
	v := validator.New()

	if err := v.RegisterValidation("hasCommonWord", ContainsCommonWord); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		data TestCommonPasswordStruct
		want bool
	}{
		{
			name: "Contains Common Word ok",
			data: TestCommonPasswordStruct{CommonPassword: "kfFF3@k"},
			want: true,
		},
		{
			name: "Contains Common Word qwerty123 fail",
			data: TestCommonPasswordStruct{CommonPassword: "Ffgf!qwerty123"},
			want: false,
		},
		{
			name: "Contains Common Word 12345678 fail",
			data: TestCommonPasswordStruct{CommonPassword: "Ffgf!12345678"},
			want: false,
		},
		{
			name: "Contains Common Word 87654321 fail",
			data: TestCommonPasswordStruct{CommonPassword: "Ffgf!87654321"},
			want: false,
		},
		{
			name: "Contains Common Word password fail",
			data: TestCommonPasswordStruct{CommonPassword: "password4321"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Struct(tt.data); (err == nil) != tt.want {
				t.Errorf("InThePast() = %v, want %v", (err == nil), tt.want)
			}
		})
	}
}
func TestNoContainsSequence(t *testing.T) {
	v := validator.New()

	if err := v.RegisterValidation("hasSequence", NoContainsSequence); err != nil {
		t.Error(err)
	}
	tests := []struct {
		name string
		data TestSequenceStruct
		want bool
	}{
		{
			name: "Sequence Password ok",
			data: TestSequenceStruct{SequencePassword: "dknfkglnfk!"},
			want: true,
		},
		{
			name: "Sequence Password fail 123",
			data: TestSequenceStruct{SequencePassword: "dknfkglnfk!123"},
			want: false,
		},
		{
			name: "Sequence Password abc fail",
			data: TestSequenceStruct{SequencePassword: "abcdknfkglnfk!"},
			want: false,
		},
		{
			name: "Sequence Password xyz fail ",
			data: TestSequenceStruct{SequencePassword: "abcdknfkglnfkxyz!"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Struct(tt.data); (err == nil) != tt.want {
				t.Errorf("InThePast() = %v, want %v", (err == nil), tt.want)
			}
		})
	}
}

func TestContainsDigit(t *testing.T) {
	v := validator.New()

	if err := v.RegisterValidation("hasDigit", ContainsDigit); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		data TestDigitStruct
		want bool
	}{
		{
			name: "Contains digit ok",
			data: TestDigitStruct{DigitPassword: "dknfkglnfk!1"},
			want: true,
		},
		{
			name: "Contains digit fail",
			data: TestDigitStruct{DigitPassword: "dknfkglnfk!"},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Struct(tt.data); (err == nil) != tt.want {
				t.Errorf("InThePast() = %v, want %v", (err == nil), tt.want)
			}
		})
	}
}
func TestContainsSpecialChar(t *testing.T) {
	v := validator.New()

	if err := v.RegisterValidation("hasSpecialChar", ContainsSpecialChar); err != nil {
		t.Error(err)
	}

	tests := []struct {
		name string
		data TestSpecialCharStruct
		want bool
	}{
		{
			name: "Special Char Password ok",
			data: TestSpecialCharStruct{SpecialCharPassword: "dsfwefFFF!F135"},
			want: true,
		},
		{
			name: "Special Char Password fail",
			data: TestSpecialCharStruct{SpecialCharPassword: "fjdhkjh387"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := v.Struct(tt.data); (err == nil) != tt.want {
				t.Errorf("InThePast() = %v, want %v", (err == nil), tt.want)
			}
		})
	}

}

func TestContainsUpper(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"valid", "A", true},
		{"invalid", "a", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ContainsUpper(testField{v: tt.val}); got != tt.want {
				t.Errorf("ContainsUpper() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testField struct {
	v string
	validator.FieldLevel
}

func (t testField) Field() reflect.Value {
	return reflect.ValueOf(t.v)
}
