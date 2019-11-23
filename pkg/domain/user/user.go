package user

import (
	"github.com/steams/salmon-hub/pkg/storage"
)

type User struct {
	Username string
	Password string
}

func Add(username, password string) {
	storage.UserStore().Add(username, password)
}

func Login(u User) string {
	return storage.UserStore().Get(u.Username, u.Password).Id
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
