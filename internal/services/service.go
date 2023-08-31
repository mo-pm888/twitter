package services

import "net/http"

//go:generate mockgen -source=service.go -destination=mocks/mock.go
type Authorization interface {
	AuthHandler(next http.Handler) http.Handler
	AdminAuthHandler(next http.Handler) http.Handler
}

type Tweets interface {
	Create(w http.ResponseWriter, r *http.Request)
	DeleteTweet(w http.ResponseWriter, r *http.Request)
	Edit(w http.ResponseWriter, r *http.Request)
	Home(w http.ResponseWriter, r *http.Request)
}

type Users interface {
	Create(w http.ResponseWriter, r *http.Request)
	EditProfile(w http.ResponseWriter, r *http.Request)
	Follow(w http.ResponseWriter, r *http.Request)
	Unfollow(w http.ResponseWriter, r *http.Request)
	GetCurrentProfile(w http.ResponseWriter, r *http.Request)
	GetAllFollowers(w http.ResponseWriter, r *http.Request)
	LogIn(w http.ResponseWriter, r *http.Request)
	LogOut(w http.ResponseWriter, r *http.Request)
}
