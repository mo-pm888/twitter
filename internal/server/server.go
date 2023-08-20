package server

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"Twitter_like_application/config"
	"Twitter_like_application/internal/admin"
	Tweets "Twitter_like_application/internal/tweets"
	Serviceuser "Twitter_like_application/internal/users"

	"github.com/gorilla/mux"
)

func Server(c config.Config) error {
	r := mux.NewRouter()
	fmt.Printf("starting server on %s:%s", c.ServerHost, c.ServerPort)
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
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.Stats)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id_user}/followers", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetAllFollowers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id_user}/followings", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetAllFollowings)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/{id_user}/block", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.BlockUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/users/{id_user}/unblock", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.UnblockUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/users/get_unblock", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.GetAllUnblockUsers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/users/get_block", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.GetAllBlockUsers)).ServeHTTP(w, r)
	}).Methods(http.MethodGet)
	r.HandleFunc("/v1/tweets/{id_tweet}/block", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.BlockTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets/{id_tweet}/unblock", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AdminAuthHandler(http.HandlerFunc(admin.UnblockTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
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

		recorder.Body.WriteTo(w)
	})
}
func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}
