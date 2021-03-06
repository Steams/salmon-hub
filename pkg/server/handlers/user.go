package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/session"
	"github.com/steams/salmon-hub/pkg/user"
	"net/http"
	"os"
)

func Login_handler(user_service user.Service, session_service session.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "POST":
			var u user.User

			fmt.Println("POST LOGIN")
			if r.Body == nil {
				http.Error(w, "Please send a request body", 400)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&u)
			fmt.Fprintln(os.Stdout, "login: %+v", u)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user_id, err := user_service.Login(u.Username, u.Password)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			var session_token string
			var csrf_token string

			session_token, err = session_service.Retrieve(user_id)
			if err != nil {
				session_token, csrf_token = session_service.Create(user_id)
			} else {
				csrf_token = session_service.Csrf(session_token)
			}

			cookie := http.Cookie{Name: "session_token", Value: session_token, HttpOnly: true, Path: "/"}
			http.SetCookie(w, &cookie)

			w.WriteHeader(http.StatusOK)

			js, err := json.Marshal(csrf_token)
			w.Write(js)

		}

	}

}

func Csrf_handler(user_service user.Service, session_service session.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "GET":

			cookie, err := r.Cookie("session_token")

			if err != nil {
				http.Error(w, "session token not present", http.StatusBadRequest)
				return
			}

			_, err = session_service.Resolve(cookie.Value)
			if err != nil {
				http.Error(w, "session expired", http.StatusBadRequest)
				return
			}

			csrf_token := session_service.Csrf(cookie.Value)

			js, err := json.Marshal(csrf_token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}

	}

}

func Register_handler(user_service user.Service, session_service session.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "POST":
			var u user.User

			fmt.Println("POST LOGIN")
			if r.Body == nil {
				http.Error(w, "Please send a request body", 400)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&u)
			fmt.Fprintln(os.Stdout, "login: %+v", u)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			user_id, err := user_service.Login(u.Username, u.Password)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusOK)
			js, err := json.Marshal(user_id)
			w.Write(js)
		}

	}

}

func Signup_handler(service user.Service) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// enableCors(&w) NOTE You dont want people signing up over the raw api

		switch r.Method {
		case "POST":
			var u user.SignupForm

			fmt.Println("POST SIGNUP")
			if r.Body == nil {
				http.Error(w, "Please send a request body", 400)
				return
			}

			err := json.NewDecoder(r.Body).Decode(&u)
			fmt.Fprintln(os.Stdout, "login: %+v", u)

			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// token := user.Signup(u.Username, u.Password, u.Email)
			service.Signup(u.Username, u.Password, u.Email)

			res, err := service.Login(u.Username, u.Password)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			cookie := http.Cookie{Name: "session_token", Value: res, HttpOnly: true}

			http.SetCookie(w, &cookie)

			w.WriteHeader(http.StatusOK)
		}

	}
}

// func signup_handler(w http.ResponseWriter, r *http.Request) {
// 	// enableCors(&w) NOTE You dont want people signing up over the raw api

// 	switch r.Method {
// 	case "POST":
// 		var u user.SignupForm

// 		fmt.Println("POST SIGNUP")
// 		if r.Body == nil {
// 			http.Error(w, "Please send a request body", 400)
// 			return
// 		}

// 		err := json.NewDecoder(r.Body).Decode(&u)
// 		fmt.Fprintln(os.Stdout, "login: %+v", u)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		token := user.Signup(u.Username, u.Password, u.Email)
// 		// js, err := json.Marshal(user.Login(u.Username, u.Password))
// 		js, err := json.Marshal(token)
// 		w.Write(js)
// 	}

// }

// func verification_handler(w http.ResponseWriter, r *http.Request) {

// 	// enableCors(&w) NOTE You dont want people signing up over the raw api

// 	switch r.Method {
// 	case "POST":
// 		var token string

// 		fmt.Println("POST SIGNUP")
// 		if r.Body == nil {
// 			http.Error(w, "Please send a request body", 400)
// 			return
// 		}

// 		err := json.NewDecoder(r.Body).Decode(&token)
// 		fmt.Fprintln(os.Stdout, "login: %+v", token)

// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusBadRequest)
// 			return
// 		}

// 		form, error := user.Verify(token)
// 		if error != nil {
// 			http.Error(w, error.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		user.Add(form.Username, form.Password, form.Email)
// 		js, err := json.Marshal(user.Login(form.Username, form.Password))
// 		w.Write(js)
// 	}

// }
