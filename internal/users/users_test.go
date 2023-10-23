package users

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"golang.org/x/crypto/bcrypt"
)

var jsonResponse struct {
	ErrText string `json:"errtext"`
}

func TestCreateCheckPwd(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()
	recorder := httptest.NewRecorder()

	service := &Service{DB: mockDB}

	newUser := &User{
		Name:      "test_name",
		Email:     "test@example.com",
		Password:  "TestPassword123",
		BirthDate: "1988-08-08",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Password = regexp.QuoteMeta(string(hashedPassword))
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	request, err := http.NewRequest("POST", "/v1/users/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	queryInsert := regexp.QuoteMeta(`INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	mock.ExpectQuery(queryInsert).
		WithArgs(newUser.Name, newUser.Password, newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).
		WillReturnError(sql.ErrNoRows)

	service.Create(recorder, request)

	err = json.NewDecoder(strings.NewReader(recorder.Body.String())).Decode(&jsonResponse)
	if err != nil {
		return
	}
	re := regexp.MustCompile(`\[string - ([^]]+)\]`)

	matches := re.FindAllStringSubmatch(jsonResponse.ErrText, -1)

	if len(matches) > 0 && len(matches[0]) > 1 {
		value := matches[0][1]
		newValue := strings.ReplaceAll(value, "\\", "")

		hashedPasswordFromDB := newValue

		userEnteredPassword := "TestPassword123"

		err := bcrypt.CompareHashAndPassword([]byte(hashedPasswordFromDB), []byte(userEnteredPassword))

		if err == nil {
		}
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println("password's wronged")
		return
	} else {
		fmt.Println("compare wrong:", err)
		return
	}

}

func TestCreateOK(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()
	recorder := httptest.NewRecorder()

	service := &Service{DB: mockDB}

	newUser := &User{
		Name:      "test_name",
		Email:     "test@example.com",
		Password:  "TestPassword123",
		BirthDate: "1988-08-08",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Password = regexp.QuoteMeta(string(hashedPassword))
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	request, err := http.NewRequest("POST", "/v1/users/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	query := regexp.QuoteMeta(`SELECT id FROM users_tweeter WHERE email = $1`)
	mock.ExpectQuery(query).
		WithArgs(newUser.Email).
		WillReturnError(sql.ErrNoRows)

	queryInsert := regexp.QuoteMeta(`INSERT INTO users_tweeter (name, password, email, nickname, location, bio, birthdate) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`)
	mock.ExpectQuery(queryInsert).
		WithArgs(newUser.Name, sqlmock.AnyArg(), newUser.Email, newUser.Nickname, newUser.Location, newUser.Bio, newUser.BirthDate).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	service.Create(recorder, request)

	assert.Equal(t, http.StatusCreated, recorder.Code)

	expectedResponse := ` "new user was created" `
	actualResponse := recorder.Body.String()

	if strings.TrimSpace(expectedResponse) != strings.TrimSpace(actualResponse) {
		t.Errorf("expected response: %s, got: %s", expectedResponse, actualResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}

func TestCreateThisUserIs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()
	recorder := httptest.NewRecorder()

	service := &Service{DB: mockDB}

	newUser := &User{
		Name:      "test_name",
		Email:     "test@example.com",
		Password:  "TestPassword123",
		BirthDate: "1988-08-08",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Password = regexp.QuoteMeta(string(hashedPassword))
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	request, err := http.NewRequest("POST", "/v1/users/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	query := regexp.QuoteMeta(`SELECT id FROM users_tweeter WHERE email = $1`)
	mock.ExpectQuery(query).
		WithArgs(newUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	service.Create(recorder, request)
	expectedResponse := ` {"errtext":"The user has already existed with this email"}`
	actualResponse := recorder.Body.String()

	if strings.TrimSpace(expectedResponse) != strings.TrimSpace(actualResponse) {
		t.Errorf("expected response: %s, got: %s", expectedResponse, actualResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
func TestCreateErrSQL(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer mockDB.Close()
	recorder := httptest.NewRecorder()

	service := &Service{DB: mockDB}

	newUser := &User{
		Name:      "test_name",
		Email:     "test@example.com",
		Password:  "TestPassword123",
		BirthDate: "1988-08-08",
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}
	newUser.Password = regexp.QuoteMeta(string(hashedPassword))
	requestBody, err := json.Marshal(newUser)
	if err != nil {
		t.Fatalf("failed to marshal request body: %v", err)
	}
	request, err := http.NewRequest("POST", "/v1/users/create", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	query := regexp.QuoteMeta(`SELECT id FROM users_tweeter WHERE email = $1`)
	mock.ExpectQuery(query).
		WithArgs(newUser.Email).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(nil))

	service.Create(recorder, request)

	expectedResponse := `{"errtext":"sql: Scan error on column index 0, name \"id\": converting NULL to int is unsupported"}`
	actualResponse := recorder.Body.String()

	if strings.TrimSpace(expectedResponse) != strings.TrimSpace(actualResponse) {
		t.Errorf("expected response: %s, got: %s", expectedResponse, actualResponse)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled expectations: %v", err)
	}
}
