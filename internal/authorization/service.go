package authorization

import (
	"context"
	"database/sql"
	"net/http"

	"Twitter_like_application/internal/services"
)

const (
	ctxKeyUserID  = "userID"
	ctxKeyIsAdmin = "isAdmin"
)

type Service struct {
	DB *sql.DB
}

func New(db *sql.DB) *Service {
	return &Service{DB: db}
}

func checkAuth(w http.ResponseWriter, r *http.Request, s *sql.DB) *http.Request {
	apikey := r.Header.Get("X-API-KEY")
	cookie, err := r.Cookie("session")
	if apikey == "" && (err != nil || cookie == nil) {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	var sessionID string
	if apikey != "" {
		sessionID = apikey
	} else if cookie != nil {
		sessionID = cookie.Value
	}
	if sessionID == "" {
		services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		return nil
	}
	query := `SELECT us.user_id, ut.admin FROM user_session us JOIN users_tweeter ut ON us.user_id = ut.id WHERE us.session_id = $1 LIMIT 1`

	var userID string
	var isAdmin bool
	err = s.QueryRow(query, sessionID).Scan(&userID, &isAdmin)
	if err != nil {
		services.ReturnErr(w, err.Error(), http.StatusInternalServerError)
		return nil
	}
	ctx := context.WithValue(r.Context(), ctxKeyUserID, userID)
	ctx = context.WithValue(ctx, ctxKeyIsAdmin, isAdmin)
	return r.WithContext(ctx)
}

func (s *Service) AuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = checkAuth(w, r, s.DB)
		if v, ok := r.Context().Value(ctxKeyUserID).(string); ok && v != "" {
			next.ServeHTTP(w, r)
		} else {
			services.ReturnErr(w, "Unauthorized", http.StatusUnauthorized)
		}

	})
}

func (s *Service) AdminAuthHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r = checkAuth(w, r, s.DB)
		if v, ok := r.Context().Value(ctxKeyIsAdmin).(bool); ok && v {
			next.ServeHTTP(w, r)
		} else {
			services.ReturnErr(w, "Unauthorized as an admin", http.StatusUnauthorized)
		}
	})
}
