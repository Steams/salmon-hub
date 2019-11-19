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
    playlist text,
    hash text
);
`

type Media struct {
	Title    string
	Duration int
	Playlist string
	Hash     string
}

type MediaStorage interface {
	Add(m Media)
	List() []Media
	ListHashes() []string
	Delete(hash string)
}

type media_storage_impl struct {
	db *sqlx.DB
}

func (x media_storage_impl) Delete(hash string) {
	rows, err := x.db.Queryx("DELETE FROM Media WHERE hash = $1", hash)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
	}
}

func (x media_storage_impl) Add(m Media) {
	stmt, err := x.db.Prepare("INSERT INTO Media(title, duration, playlist, hash) values(?,?,?,?)")
	if err != nil {
		panic(err)
	}
	_, _ = stmt.Exec(m.Title, m.Duration, m.Playlist, m.Hash)
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

	if media == nil {
		media = make([]Media, 0)
		return media
	}

	return media
}

func (x media_storage_impl) ListHashes() []string {

	rows, err := x.db.Queryx("Select hash FROM media")
	if err != nil {
		log.Fatalln(err)
	}

	hash := ""
	var hashes []string

	for rows.Next() {
		err := rows.Scan(&hash)
		if err != nil {
			log.Fatalln(err)
		}
		hashes = append(hashes, hash)
	}

	if hashes == nil {
		hashes = make([]string, 0)
		return hashes
	}

	return hashes
}

func Open() MediaStorage {
	os.Remove("./salmon.db")

	db, err := sqlx.Open("sqlite3", "./salmon.db")
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(schema)
	_, _ = stmt.Exec()

	return media_storage_impl{db}
}
