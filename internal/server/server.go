package server

import (
	Tweets "Twitter_like_application/internal/tweets"
	Serviceuser "Twitter_like_application/internal/users"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func Server() {
	r := mux.NewRouter()
	fmt.Println("Server was run", "localhost:8080")
	http.ListenAndServe("localhost:8080", r)
	r.HandleFunc("/v1/users", Serviceuser.CreateUser).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/login", Serviceuser.LoginUsers).Methods(http.MethodPost)
	http.Handle("/v1/users/logout", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.LogoutUser)))
	http.Handle("/v1/users/{id}", Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.GetCurrentProfile)))
	r.HandleFunc("/v1/users/reset-password", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.ResetPassword)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/follow", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.FollowUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users/unfollow", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.UnfollowUser)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/users", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Serviceuser.EditProfile)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)
	r.HandleFunc("/v1/tweets", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.CreateTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPost)
	r.HandleFunc("/v1/tweets", func(w http.ResponseWriter, r *http.Request) {
		Serviceuser.AuthHandler(http.HandlerFunc(Tweets.EditTweet)).ServeHTTP(w, r)
	}).Methods(http.MethodPatch)

}
