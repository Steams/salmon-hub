package user

import (
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
