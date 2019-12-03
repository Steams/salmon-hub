package server

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/domain/media"
	"github.com/steams/salmon-hub/pkg/domain/user"
	"log"
	"net/http"
	"os"
)

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func verification_handler(w http.ResponseWriter, r *http.Request) {

	// enableCors(&w) NOTE You dont want people signing up over the raw api

	switch r.Method {
	case "POST":
		var token string

		fmt.Println("POST SIGNUP")
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&token)
		fmt.Fprintln(os.Stdout, "login: %+v", token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		form, error := user.Verify(token)
		if error != nil {
			http.Error(w, error.Error(), http.StatusUnauthorized)
			return
		}

		user.Add(form.Username, form.Password, form.Email)
		js, err := json.Marshal(user.Login(form.Username, form.Password))
		w.Write(js)
	}

}

func signup_handler_stub(w http.ResponseWriter, r *http.Request) {
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
		user.Signup(u.Username, u.Password, u.Email)

		res := user.Login(u.Username, u.Password)
		if res == "" {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		cookie := http.Cookie{Name: "session_token", Value: res, HttpOnly: true}

		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusOK)
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

func login_handler(w http.ResponseWriter, r *http.Request) {

	defer enableCors(&w)

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

		res := user.Login(u.Username, u.Password)
		if res == "" {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		cookie := http.Cookie{Name: "session_token", Value: res, HttpOnly: true}

		http.SetCookie(w, &cookie)

		w.WriteHeader(http.StatusOK)
	}

}

func synch_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	switch r.Method {
	case "GET":

		id, err := r.Cookie("session_token")

		if err != nil {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		hashes := media.ListHashes(id.Value)

		js, err := json.Marshal(hashes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	case "POST":
		var hashes []string

		id, err := r.Cookie("session_token")

		if err != nil {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		fmt.Println("POST DELETE REQUEST")
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err = json.NewDecoder(r.Body).Decode(&hashes)
		fmt.Fprintln(os.Stdout, "HASHES: %+v", hashes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, hash := range hashes {
			media.Delete(id.Value, hash)
		}
	}

}

func media_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	switch r.Method {
	case "GET":

		id, err := r.Cookie("session_token")

		if err != nil {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		media_list := media.List(id.Value)

		if media_list == nil {
			js, _ := json.Marshal([]media.Media{})
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return
		}

		js, err := json.Marshal(media_list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return

	case "POST":

		id, err := r.Cookie("session_token")

		if err != nil {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		// NOTE Perhaps its idiomatic in go to call this media.Model instead ?
		var songs []media.Media

		fmt.Println("POST REQUEST")
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err = json.NewDecoder(r.Body).Decode(&songs)
		fmt.Fprintln(os.Stdout, "BODY: \n%+v", songs)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, song := range songs {
			fmt.Println(song)
			media.Add(id.Value, song)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println(r.URL.Path)
	if r.URL.Path == "/elm.js" {
		http.ServeFile(w, r, "/home/steams/Development/audigo/salmon-web-client/elm.js")
		return
	}
	// path := "/home/steams/Development/audigo/salmon-web-client/" + r.URL.Path[1:]
	path := "/home/steams/Development/audigo/salmon-web-client/index.html"
	http.ServeFile(w, r, path)
}

func Run() {

	user.Add("admin", "password", "email")
	fmt.Println(user.Get("admin", "password"))

	http.HandleFunc("/", file_handler)
	http.HandleFunc("/media", media_handler)
	http.HandleFunc("/synch", synch_handler)
	http.HandleFunc("/api/login", login_handler)
	http.HandleFunc("/api/signup", signup_handler_stub)
	http.HandleFunc("/api/verify", verification_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
