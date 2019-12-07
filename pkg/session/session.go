package session

import (
	"github.com/steams/salmon-hub/pkg/rand"
	"strings"
)

const (
	SECRET_KEY = "secret"
)

type Service interface {
	Create(user_id string) string
	Get(session_token string) string
	Delete(session_token string)
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

func (s service_imp) Create(user_id string) string {
	id := generate_large_random()

	s.repo.Add(user_id, id)

	token := to_token(id)
	return token

}

func (s service_imp) Get(session_token string) string {
	id := from_token(session_token)

	user_id := s.repo.Get(id)
	return user_id

}

func (s service_imp) Delete(session_token string) {
	id := from_token(session_token)

	s.repo.Delete(id)
}

func to_token(session_id string) string {
	return encrypt("SESSION-TOKEN:"+session_id, SECRET_KEY)
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
