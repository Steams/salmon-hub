package storage

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

const media_schema = `
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
const user_schema = `
CREATE TABLE user (
    id text,
    username text,
    password text
);
`

var db *sqlx.DB

func init() {
	fmt.Println("Initializing db")
	os.Remove("./salmon.db")

	var err error
	db, err = sqlx.Open("sqlite3", "./salmon.db")

	if err != nil {
		panic(err)
	}

	db.MustExec(media_schema)
	db.MustExec(user_schema)
}

func MediaStore() MediaStorage {
	return media_storage_impl{db}
}

func UserStore() UserStorage {
	return user_storage_impl{db}
}
