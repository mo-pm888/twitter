package services

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReturnErr(t *testing.T) {
	type args struct {
		w    http.ResponseWriter
		err  string
		code int
	}
	err := errors.New("sample error")
	errString := err.Error()
	tests := []struct {
		name         string
		args         args
		expectedJSON string
	}{
		{
			name: "Test with string error",
			args: args{
				w:    httptest.NewRecorder(),
				err:  "sample error message",
				code: http.StatusInternalServerError,
			},
			expectedJSON: `{"errtext":"sample error message"}
`,
		},
		{
			name: "Test with error",
			args: args{
				w:    httptest.NewRecorder(),
				err:  errString,
				code: http.StatusBadRequest,
			},
			expectedJSON: `{"errtext":"sample error"}
`,
		},
		{
			name: "Test with empty error",
			args: args{
				w:    httptest.NewRecorder(),
				err:  "",
				code: http.StatusNotFound,
			},
			expectedJSON: `{"errtext":""}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReturnErr(tt.args.w, tt.args.err, tt.args.code)

			if tt.args.w.(*httptest.ResponseRecorder).Code != tt.args.code {
				t.Errorf("Expected status code %d, but got %d", tt.args.code, tt.args.w.(*httptest.ResponseRecorder).Code)
			}

			contentType := tt.args.w.(*httptest.ResponseRecorder).Header().Get("Content-Type")
			if contentType != "application/json; charset=utf-8" {
				t.Errorf("Expected Content-Type %s, but got %s", "application/json; charset=utf-8", contentType)
			}

			responseBody := tt.args.w.(*httptest.ResponseRecorder).Body.String()
			if responseBody != tt.expectedJSON {
				t.Error(responseBody, tt.expectedJSON)
				t.Errorf("Expected response body:%s but got:%s", tt.expectedJSON, responseBody)
			}

		})
	}
}
