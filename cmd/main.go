package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/steams/salmon-hub/pkg/media"
	media_repo "github.com/steams/salmon-hub/pkg/media/sqlite_repo"
	"github.com/steams/salmon-hub/pkg/server"
	"github.com/steams/salmon-hub/pkg/user"
	user_repo "github.com/steams/salmon-hub/pkg/user/sqlite_repo"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {

	fmt.Println("Initializing db")
	os.Remove("./salmon.db")

	db, err := sqlx.Open("sqlite3", "./salmon.db")

	if err != nil {
		panic(err)
	}

	db.MustExec(user_repo.Schema)
	db.MustExec(media_repo.Schema)

	// since this is an sqlite repo, you;ll never test it separate from an sqlite db, so u might aswell pass in the db instead of the store interface and create the interface inside, there is no benifit from being able to pass in a different "UserStore" since the repo specifies it must be an sqlite one anyways
	userRepo := user_repo.New(db)
	userService := user.CreateService(userRepo)

	userService.Signup("admin", "password", "email")

	mediaRepo := media_repo.New(db)
	mediaService := media.CreateService(mediaRepo)

	server := server.New(userService, mediaService, "8080")

	if err = server.Run(); err != nil {
		return err
	}

	return nil
}
