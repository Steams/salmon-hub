package storage

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type Media struct {
	Title    string
	Artist   string
	Album    string
	Duration int
	Playlist string
	Hash     string
	UserID   string
}

type MediaStorage interface {
	Add(m Media)
	List(user_id string) []Media
	ListHashes(user_id string) []string
	Delete(user_id, hash string)
}

type media_storage_impl struct {
	db *sqlx.DB
}

func (x media_storage_impl) Delete(user_id, hash string) {
	rows, err := x.db.Queryx("DELETE FROM Media WHERE hash = $1 and userid =$2", hash, user_id)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
	}
}

func (x media_storage_impl) Add(m Media) {
	stmt, err := x.db.Preparex("INSERT INTO Media(title, artist, album, duration, playlist, hash, userid) values(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(m.Title, m.Artist, m.Album, m.Duration, m.Playlist, m.Hash, m.UserID)
}

func (x media_storage_impl) List(user_id string) []Media {

	rows, err := x.db.Queryx("Select * FROM media WHERE userid=$1", user_id)
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

func (x media_storage_impl) ListHashes(user_id string) []string {

	rows, err := x.db.Queryx("Select hash FROM media WHERE userid=$1", user_id)
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
