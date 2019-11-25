package user

import (
	"errors"
	"github.com/steams/salmon-hub/pkg/storage"
)

type User struct {
	Username string
	Password string
}

type SignupForm struct {
	Username string
	Password string
	Email    string
}

// take user info,generate a jwt token containing user info and use it as a link, when the link is sent back to verify the token, extract info
// if link is used twice, just check if the user already exists
// you can include an expirable token in the jwt if you want

// func Signup(username, password, email string) string {
// 	// return (username + "-" + password + "-" + email)
// }
func Signup(username, password, email string) {
	Add(username, password, email)
}

func Verify(jwt string) (SignupForm, error) {
	//extract user info from jwt
	if jwt != "admin-password-email" {
		return SignupForm{}, errors.New("invalid verification token")
	}
	return SignupForm{"admin", "password", "email"}, nil
}

func Add(username, password, email string) {
	storage.UserStore().Add(username, password)
}

func Login(username, password string) string {
	return storage.UserStore().Get(username, password).Id
}

func Get(username, password string) User {
	return from_db(storage.UserStore().Get(username, password))
}

func from_db(u storage.User) User {
	return User{
		Username: u.Username,
		Password: u.Password,
	}
}
