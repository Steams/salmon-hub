package storage

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

type User struct {
	Id       string
	Username string
	Password string
}

type UserStorage interface {
	Add(username, password string)
	Get(username, password string) User
	Delete(id string)
}

type user_storage_impl struct {
	db *sqlx.DB
}

func (x user_storage_impl) Delete(id string) {
	rows, err := x.db.Queryx("DELETE FROM User WHERE id = $1", id)
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
	}
}

func (x user_storage_impl) Add(username, password string) {
	stmt, err := x.db.Preparex("INSERT INTO User(id, username, password ) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(uuid.New().String(), username, password)
	// stmt.MustExec("1", username, password)
}

func (x user_storage_impl) Get(username, password string) User {

	rows, err := x.db.Queryx("Select * FROM User WHERE username = $1 AND password = $2", username, password)
	if err != nil {
		log.Fatalln(err)
	}

	u := User{}

	for rows.Next() {
		err := rows.StructScan(&u)
		if err != nil {
			log.Fatalln(err)
		}
		break
	}
	rows.Close()

	return u
}
