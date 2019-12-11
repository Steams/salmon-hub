package session

import (
	// "crypto/aes"
	// "crypto/cipher"
	// crypto_rand "crypto/rand"
	"errors"
	// "fmt"
	"github.com/steams/salmon-hub/pkg/rand"
	// "io"
	"strings"
)

const (
	SECRET_KEY = "16bytesecrethere"
)

type Service interface {
	Create(user_id string) (string, string)
	Retrieve(user_id string) (string, error)
	Resolve(session_token string) (string, error)
	Delete(session_token string)
	Csrf(session_token string) string
	Validate(session_token, csrf_token string) bool
}

type Repository interface {
	Add(user_id, session_id string)
	Resolve(session_id string) string
	Delete(session_id string)
	Retrieve(user_id string) string
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

func (s service_imp) Retrieve(user_id string) (string, error) {
	session_id := s.repo.Retrieve(user_id)
	if session_id == "" {
		return "", errors.New("no session")
	}

	return to_token(session_id), nil
}

func (s service_imp) Resolve(session_token string) (string, error) {
	id := from_token(session_token)

	user_id := s.repo.Resolve(id)
	if user_id == "" {
		return "", errors.New("no session")
	}
	return user_id, nil
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

func from_token(token string) string {
	//decrypt token,split token on colon and take second value in array
	return strings.Split(decrypt(token, SECRET_KEY), ":")[1]
}

func to_csrf(session_id string) string {
	return encrypt("CSRF-TOKEN:"+session_id, SECRET_KEY)
}

func generate_large_random() string {
	return rand.Hex_128()
}

// TODO These should be in another module
func encrypt(s, secret string) string {
	// key := []byte(secret)

	// plaintext := []byte(s)

	// block, err := aes.NewCipher(key)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// nonce := make([]byte, 12)
	// if _, err := io.ReadFull(crypto_rand.Reader, nonce); err != nil {
	// 	panic(err.Error())
	// }

	// aesgcm, err := cipher.NewGCM(block)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)
	// return string(ciphertext)
	return s

}

func decrypt(s, secret string) string {
	// key := []byte(secret)
	// ciphertext := []byte(s)

	// block, err := aes.NewCipher(key)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// aesgcm, err := cipher.NewGCM(block)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// nonceSize := aesgcm.NonceSize()
	// if len(ciphertext) < nonceSize {
	// 	fmt.Println(err)
	// }

	// nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Printf("%s\n", plaintext)
	// return string(plaintext)
	return s
}
