package session

import (
	"github.com/steams/salmon-hub/pkg/rand"
	"strings"
)

const (
	SECRET_KEY = "secret"
)

type Service interface {
	Create(user_id string) (string, string)
	Get(session_token string) string
	Delete(session_token string)
	Csrf(session_token string) string
	Validate(session_token, csrf_token string) bool
}

type Repository interface {
	Add(user_id, session_id string)
	Get(session_id string) string
	Delete(session_id string)
}

type service_imp struct {
	repo Repository
}

func CreateService(repo Repository) Service {
	return service_imp{repo}
}

func (s service_imp) Create(user_id string) (session_token string, csrf_token string) {
	id := generate_large_random()

	s.repo.Add(user_id, id)

	token := to_token(id)
	csrfToken := to_csrf(id)

	return token, csrfToken
}

func (s service_imp) Get(session_token string) string {
	id := from_token(session_token)

	user_id := s.repo.Get(id)
	return user_id

}

// TODO this needs to vaidatte that the id actually existing
func (s service_imp) Csrf(session_token string) string {
	return to_csrf(from_token(session_token))
}

func (s service_imp) Delete(session_token string) {
	id := from_token(session_token)

	s.repo.Delete(id)
}

func (s service_imp) Validate(session_token, csrf_token string) bool {
	return from_token(session_token) == from_token(csrf_token)
}

func to_token(session_id string) string {
	return encrypt("SESSION-TOKEN:"+session_id, SECRET_KEY)
}

func to_csrf(session_id string) string {
	return encrypt("CSRF-TOKEN:"+session_id, SECRET_KEY)
}

func from_token(token string) string {
	//decrypt token,split token on colon and take second value in array
	return strings.Split(decrypt(token, SECRET_KEY), ":")[1]
}

func generate_large_random() string {
	return rand.Hex_128()
}

func encrypt(s, secret string) string {
	return s
}

func decrypt(s, secret string) string {
	return s
}
