package sqlite_repo

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/steams/salmon-hub/pkg/media"
	"log"
)

type Repository = media.Repository
type Media = media.Media

const Schema = `
CREATE TABLE media (
    title text,
    artist text,
    album text,
    duration int,
    playlist text,
    hash text,
    userid text
)
`

func New(db *sqlx.DB) Repository {
	return repository{db}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) Add(id string, m Media) {

	stmt, err := r.db.Preparex("INSERT INTO Media(title, artist, album, duration, playlist, hash, userid) values(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(m.Title, m.Artist, m.Album, m.Duration, m.Playlist, m.Hash, id)

}

func (r repository) List(user_id string) []Media {

	rows, err := r.db.Queryx("Select title,artist,album,duration,playlist,hash FROM media WHERE userid=$1", user_id)
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

func (r repository) ListHashes(user_id string) []string {

	rows, err := r.db.Queryx("Select hash FROM media WHERE userid=$1", user_id)
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

func (r repository) Delete(user_id string, hash string) {

	stmt, err := r.db.Preparex("DELETE FROM Media WHERE hash = ? and userid = ?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(hash, user_id)
}
