package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/admin"
	tweets "Twitter_like_application/internal/tweets"
	"Twitter_like_application/internal/users"

	"github.com/gorilla/mux"
)

func Server(c config.Config, s users.Service, t tweets.Service, a admin.Service) error {
	r := mux.NewRouter()
	fmt.Printf("starting server on %s:%s\n", c.ServerHost, c.ServerPort)
	r.Use(LoggingMiddleware)
	r.Use(CorsMiddleware)
	r.HandleFunc("/v1/users/create", s.Create).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/login", s.LogIn).Methods(http.MethodPost)
	//http.Handle("/v1/users/logout", users.AuthHandler(http.HandlerFunc(s.LogOut)))
	r.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.GetCurrentProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/home", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Home)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/reset-password", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.ResetPassword)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/{id}/follow", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.Follow)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/{id}/unfollow", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.Unfollow)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/edit", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.EditProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/create", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Create)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Edit)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/{id_tweet}/retweet", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Retweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}/like", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Like)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}/unlike", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(t.Unlike)).ServeHTTP(w, r)
	}).Methods(http.MethodDelete)
	r.HandleFunc("/v1/users/{id_user}/followers", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.GetAllFollowers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id_user}/followings", func(w http.ResponseWriter, r *http.Request) {
		s.AuthHandler(http.HandlerFunc(s.GetAllFollowings)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id_user}/block", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.BlockUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/users/{id_user}/unblock", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.UnblockUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/{id_tweet}/block", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.BlockTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/{id_tweet}/unblock", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.UnblockTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/users/get_unblock", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.GetAllUnblockUsers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/get_block", func(w http.ResponseWriter, r *http.Request) {
		s.AdminAuthHandler(http.HandlerFunc(a.GetAllBlockUsers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)

	})
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", c.ServerHost, c.ServerPort), r)
	fmt.Println(err)
	return err
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s", r.Method, r.URL.Path)

		recorder := httptest.NewRecorder()

		next.ServeHTTP(recorder, r)

		log.Printf("Sent response: %d %s", recorder.Code, http.StatusText(recorder.Code))

		for k, v := range recorder.Header() {
			w.Header()[k] = v
		}
		w.WriteHeader(recorder.Code)

		_, err := recorder.Body.WriteTo(w)
		if err != nil {
			return
		}
	})
}
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
