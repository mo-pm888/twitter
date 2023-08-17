package users

import (
	Tweets "Twitter_like_application/internal/tweets"
)

type User struct {
	ID                 int
	Name               string `json:"name" validate:"omitempty"`
	Password           string `json:"password" validate:"omitempty"`
	Email              string `json:"email" validate:"omitempty,email"`
	EmailToken         string
	ConfirmEmailToken  bool
	ResetPasswordToken string
	BirthDate          string `json:"birthdate" validate:"omitempty"`
	Nickname           string `json:"nickname" validate:"omitempty"`
	Bio                string `json:"bio" validate:"omitempty"`
	Location           string `json:"location" validate:"omitempty"`
}
type GetCurrentUser struct {
	Name      string `json:"name"`
	BirthDate string `json:"birthdate"`
	Nickname  string `json:"nickname"`
	Bio       string `json:"bio"`
	Location  string `json:"location"`
	Following int    `json:"following"`
	Followers int    `json:"followers"`
}

type ReplaceMyData struct {
	NewName      string `json:"new_name"`
	NewPassword  string `json:"new_password"`
	NewEmail     string `json:"new_email"`
	NewBirthDate string `json:"new_birth_date"`
	NewNickname  string `json:"new_nickname"`
	NewBio       string `json:"new_bio"`
	NewLocation  string `json:"new_location"`
}

type ReplayTweet struct {
	Tweets.Tweet
}

type DeleteUserST struct {
	UserIdDeleting int `json:"delete_id"`
}

type ResetPasswordUser struct {
	UserResetPassword int `json:"user_reset_password"`
}

type FollowingForUser struct {
	Writer     int `json:"writer"`
	Subscriber int `json:"subscriber"`
	User
}

type UsersLogin struct {
	Usermail string `json:"email_logIN"`
	Password string `json:"password_logIN"`
}

type Tweeter_like struct {
	Autor      int `json:"autor"`
	Id_post    int `json:"id_post"`
	Whose_like int `json:"whose_like"`
}
