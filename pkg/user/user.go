package user

// contins models, defines/implements user service API, and defines repository API
import (
	"errors"
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

type Repository interface {
	Add(f SignupForm)
	Get(username, password string) string
}

type Service interface {
	Signup(username, password, email string)
	// Verify(jwt string) string
	Login(username, password string) string
}

func CreateService(r Repository) Service {
	return service_imp{r}
}

type service_imp struct {
	repo Repository
}

func (s service_imp) Signup(username, password, email string) {
	s.repo.Add(SignupForm{username, password, email})
}

func (s service_imp) Login(username, password string) string {
	return s.repo.Get(username, password)
}

// take user info,generate a jwt token containing user info and use it as a link, when the link is sent back to verify the token, extract info
// if link is used twice, just check if the user already exists
// you can include an expirable token in the jwt if you want

// func Signup(username, password, email string) string {
// 	// return (username + "-" + password + "-" + email)
// }

func Verify(jwt string) (SignupForm, error) {
	//extract user info from jwt
	if jwt != "admin-password-email" {
		return SignupForm{}, errors.New("invalid verification token")
	}
	return SignupForm{"admin", "password", "email"}, nil
}
