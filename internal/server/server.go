package server

import (
	"Twitter_like_application/internal/admin"
	Tweets "Twitter_like_application/internal/tweets"
	Serviceuser "Twitter_like_application/internal/users"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"net/http/httptest"
)

func Server() error {
	r := mux.NewRouter()
	fmt.Printf("Server was run %s:%s\n", admin.ServerHost, admin.ServerPort)
	r.Use(LoggingMiddleware)
	r.Use(CorsMiddleware)
	r.HandleFunc("/v1/users/create", Serviceuser.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/login", Serviceuser.LogIn).Methods(http.MethodPost)
	http.Handle("/v1/users/logout", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.LogOut)))
	r.HandleFunc("/v1/users/", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetCurrentProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/home", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.Home)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/reset-password", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.ResetPassword)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/{id}/follow", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.FollowUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/{id}/unfollow", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.UnfollowUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/edit", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.EditProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/create", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.CreateNewTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.EditTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/{id_tweet}/retweet", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.Retweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}/like", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.LikeTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets/{id_tweet}/unlike", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.UnlikeTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodDelete)
	r.HandleFunc("/v1/admin/stats", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(admin.Stats)).ServeHTTP(w, r)
	}).Methods(http.MethodDelete)
	r.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "*")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusOK)

	})
	err := http.ListenAndServe(fmt.Sprintf("%s:%s", admin.ServerHost, admin.ServerPort), r)
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

		recorder.Body.WriteTo(w)
	})
}
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
