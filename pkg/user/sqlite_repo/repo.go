package sqlite_repo

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/steams/salmon-hub/pkg/user"
	"log"
)

const Schema = `
CREATE TABLE user (
    id text,
    username text,
    password text
);
`

func New(db *sqlx.DB) user.Repository {
	return repository{db}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) Add(f user.SignupForm) {

	stmt, err := r.db.Preparex("INSERT INTO User(id, username, password ) values(?,?,?)")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(uuid.New().String(), f.Username, f.Password)
	// stmt.MustExec("1", username, password)
}

func (r repository) Get(username, password string) string {

	rows, err := r.db.Queryx("SELECT id FROM User WHERE username = $1 AND password = $2", username, password)
	if err != nil {
		log.Fatalln(err)
	}

	var id string = ""

	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			log.Fatalln(err)
		}
		break
	}
	rows.Close()

	return id
}

// func remove(id string) {
// 	stmt, err := x.db.Preparex("DELETE FROM User WHERE id = ?)")
// 	if err != nil {
// 		panic(err)
// 	}
// 	stmt.MustExec(id)
// }
