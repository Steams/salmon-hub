package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
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

type MediaStorage interface {
	Add(m Media)
	List() []Media
}

type media_storage_impl struct {
	db *sqlx.DB
}

func (x media_storage_impl) Add(m Media) {
	stmt, err := x.db.Prepare("INSERT INTO Media(title, duration, playlist) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	_, _ = stmt.Exec(m.Title, m.Duration, m.Playlist)
}

func (x media_storage_impl) List() []Media {

	rows, err := x.db.Queryx("Select * FROM media")
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

func Open() MediaStorage {
	os.Remove("./salmon.db")

	db, err := sqlx.Open("sqlite3", "./salmon.db")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare("CREATE TABLE Media (title varchar(255),duration int, playlist varchar(255))")
	_, _ = stmt.Exec()

	return media_storage_impl{db}
}
