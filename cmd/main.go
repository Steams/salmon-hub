package main

import (
	// "database/sql"
	"encoding/json"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
)

const schema = `
CREATE TABLE media (
    title text,
    duration int,
    playlist text
);
`

type Media struct {
	Title    string
	Duration int
	Playlist string
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

func add_media(m Media, db *sqlx.DB) {
	stmt, err := db.Prepare("INSERT INTO Media(title, duration, playlist) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	_, _ = stmt.Exec(m.Title, m.Duration, m.Playlist)
}

func get_media(db *sqlx.DB) []Media {
	rows, err := db.Queryx("Select * FROM media")
	if err != nil {
		log.Fatalln(err)
	}

	m := Media{}
	var media []Media

	for rows.Next() {
		err := rows.StructScan(&m)
		if err != nil {
			log.Fatalln(err)
		}
		media = append(media, m)
	}

	return media

}

func media_handler(db *sqlx.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":

			media_list := get_media(db)

			fmt.Fprintf(w, "Media list :\n\n%+v", media_list)

		case "POST":
			var m Media

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
			add_media(m, db)
		}
	}

}

func run() {
	os.Remove("./salmon.db")

	db, err := sqlx.Open("sqlite3", "./salmon.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("CREATE TABLE Media (title varchar(255),duration int, playlist varchar(255))")
	_, _ = stmt.Exec()

	http.HandleFunc("/media", media_handler(db))
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func main() {

	run()
}
