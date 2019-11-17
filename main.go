package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "github.com/thoas/go-funk"
	// "strings"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }

// func file_handler(w http.ResponseWriter, r *http.Request) {
// 	path := "/home/steams/Development/audigo/resources/" + r.URL.Path[1:]
// 	http.ServeFile(w, r, path)
// }

func login(username, password string) bool {
	if username != password {
		return false
	}
	return true
}

type Media struct {
	Title        string
	Duration     int
	Playlist_url string
}

type User struct {
	id       int
	name     string
	password string
	hosts    []Host
	media    []Media
}

type Host struct {
	url string
	id  int
}

// func add_media(user, media) {
// }

// func add_host(Host,store) {
// }

//return all in store where host_id == in user.hosts
// func get_media(user) {}

// func get_stream(user, media) {}

type Store interface {
	save()
	login()
}

type MockStore struct {
	media []Media
	users []User
	hosts []Host
}

func (s *MockStore) save(m Media) {
	s.media = append(s.media, m)
}

// func (s *MockStore) login(name, password string) {
// 	if name != password {
// 	}
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
// }
type server struct {
	port string
}

// func add_user(user User, db ) {
// 	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")

// }

func run() {
	hosts := []Host{
		{"localhost:8080/server", 0},
	}

	media := []Media{
		{"closer to my dreams", 180, "closer.mp3"},
		{"the comeup", 120, "comeup.mp3"},
	}

	user := User{0, "admin", "admin", hosts, media}

	db, err := sql.Open("sqlite3", "./salmon.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO user(username, password ) values(?,?)")

	res, err := stmt.Exec("admin", "adminpwd")
	if err != nil {
		panic(err)
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi there, I love %s | %s!", r.URL.Path[1:], user.name)
	}

	media_handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":

			var media_list string

			for _, x := range user.media {
				media_list = media_list + x.Title + "\n"
			}
			fmt.Fprintf(w, "Media list :\n\n%s", media_list)

		case "POST":
			var m Media

			// Try to decode the request body into the struct. If there is an error,
			// respond to the client with the error message and a 400 status code.
			fmt.Println("POST REQUEST")
			if r.Body == nil {
				http.Error(w, "Please send a request body", 400)
				return
			}
			err := json.NewDecoder(r.Body).Decode(&m)
			fmt.Fprintln(os.Stdout, "BODY: \n%+v", m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, "Media : \n%s", m.Title)
		}

		// media_list = strings.Join(funk.Map(user.media, func(x Media) string { return x.title }).([]string), "\n")

		//media_list = join "\n" $ map title user.media

	}

	http.HandleFunc("/", handler)
	http.HandleFunc("/media", media_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main() {

	run()
	// store := MockStore{}

	// store = add_media(media, store)

	// host_add_media
	// host_get_media

	//prompt for username
	//return list of media for user
	//cmd line get line
	// get_media(user)

	// if err := run(); err != nil {
	// fmt.Fprintf(os.Stderr, "%s\n", err)
	// os.Exit(1)
	// }

}
