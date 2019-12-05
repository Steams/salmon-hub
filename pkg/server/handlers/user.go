package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/user"
	"net/http"
	"os"
)

func Login_handler(service user.Service) http.HandlerFunc {
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

			res := service.Login(u.Username, u.Password)
			if res == "" {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			js, err := json.Marshal(res)
			w.Write(js)
			// cookie := http.Cookie{Name: "session_token", Value: res, HttpOnly: true}

			// http.SetCookie(w, &cookie)

			// w.WriteHeader(http.StatusOK)
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

			res := service.Login(u.Username, u.Password)
			if res == "" {
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
