package sqlite_repo

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/steams/salmon-hub/pkg/session"
	"log"
)

const Schema = `
CREATE TABLE session (
    session_id text,
    user_id text
);
`

func New(db *sqlx.DB) session.Repository {
	return repository{db}
}

type repository struct {
	db *sqlx.DB
}

func (r repository) Add(user_id, session_id string) {

	stmt, err := r.db.Preparex("INSERT INTO Session(session_id, user_id) values(?,?)")
	if err != nil {
		panic(err)
	}

	stmt.MustExec(session_id, user_id)
}

func (r repository) Get(session_id string) string {

	rows, err := r.db.Queryx("SELECT user_id FROM Session WHERE session_id = $1", session_id)

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

func (r repository) Delete(session_id string) {

	stmt, err := r.db.Preparex("DELETE FROM Session WHERE session_id = ?")
	if err != nil {
		panic(err)
	}
	stmt.MustExec(session_id)
}
