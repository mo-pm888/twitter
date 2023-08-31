//package users
//
//import (
//	"context"
//	"net/http/httptest"
//	"testing"
//
//	"github.com/DATA-DOG/go-sqlmock"
//)
//
////	func TestCreate(t *testing.T) {
////		service := new(Service)
////		jsonData := []byte(`{
////			"name": "Test User",
////			"password": "testpassword",
////	       "email": "test@example.com",
////			"birthdate": "2000-01-01"
////		}`)
////		req, err := http.NewRequest("POST", "/v1/users/create", bytes.NewBuffer(jsonData))
////		if err != nil {
////			t.Fatal(err)
////		}
////
////		recorder := httptest.NewRecorder()
////		service.Create(recorder, req)
////
////		if recorder.Code != http.StatusCreated {
////			t.Errorf("Expected %d but got %d", http.StatusCreated, recorder.Code)
////		}
////		if t.Failed() {
////			t.FailNow()
////		}
////	}
//func TestCheckAuth(t *testing.T) {
//	db, mock, err := sqlmock.New()
//	if err != nil {
//		t.Fatalf("failed to create sqlmock: %v", err)
//	}
//	defer db.Close()
//
//	r := httptest.NewRequest("GET", "/path", nil)
//	r = r.WithContext(context.WithValue(r.Context(), ctxKeyUserID, "user_id"))
//
//	query := "SELECT * FROM user_session WHERE id = ?"
//	mock.ExpectQuery(query).WithArgs("user_id").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("user_id"))
//
//	r = checkAuth(httptest.NewRecorder(), r, db)
//
//	if err := mock.ExpectationsWereMet(); err != nil {
//		t.Errorf("unfulfilled expectations: %v", err)
//	}
//}
