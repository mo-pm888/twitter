package Serviceuser

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"html/template"
	"net/http"
)

type UserDataSTR struct {
	UserData map[int]*Users
}

var (
	UserDate = make(map[int]*Users)
)

func Put(u *Users) bool {
	for _, user := range UserDate {
		if user.Email == u.Email {
			return false

		}

	}
	u.ID = len(UserDate) + 1
	UserDate[u.ID] = u

	return true

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var newUser Users
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if newUser.Name == "" || newUser.Email == "" || newUser.Password == "" || newUser.Nickname == "" {
		http.Error(w, "Invalid user data", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	newUser.Password = string(hashedPassword)
	ret := Put(&newUser)
	if ret == false {
		fmt.Fprint(w, "This user is alredy added")
		return
	} else {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
	return
}

func LoginUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		usermail := r.FormValue("usermail")
		password := r.FormValue("password")
		for _, name := range UserDate {
			if name.Email == usermail || name.Password == password {
				cookie := &http.Cookie{
					Name:  "session",
					Value: "authenticated",
				}
				http.SetCookie(w, cookie)
				http.Redirect(w, r, "/", http.StatusSeeOther)
			}
			t, _ := template.ParseFiles("login.html")
			t.Execute(w, nil)
		}
	}
}

func LogoutUser(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != http.ErrNoCookie {
		cookie = &http.Cookie{
			Name:   "session",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(w, cookie)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	var deleteUser Users
	err := json.NewDecoder(r.Body).Decode(&deleteUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for id, _ := range UserDate {
		if deleteUser.ID == id {
			delete(UserDate, id)
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

}

func Following(w http.ResponseWriter, r *http.Request) {
	var writer, subscriber Users
	err := json.NewDecoder(r.Body).Decode(&writer)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func ResetPassword(w http.ResponseWriter, r *http.Request) {

}

func Seemyprofile(w http.ResponseWriter, r *http.Request) {

}

func Home(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session")
	if err != nil || cookie.Value != "authenticated" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	t, _ := template.ParseFiles("home.html")
	t.Execute(w, nil)
}
