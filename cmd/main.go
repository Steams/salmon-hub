package main

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"github.com/steams/salmon-hub/pkg/media"
	media_repo "github.com/steams/salmon-hub/pkg/media/sqlite_repo"
	"github.com/steams/salmon-hub/pkg/server"
	"github.com/steams/salmon-hub/pkg/session"
	session_repo "github.com/steams/salmon-hub/pkg/session/sqlite_repo"
	"github.com/steams/salmon-hub/pkg/user"
	user_repo "github.com/steams/salmon-hub/pkg/user/sqlite_repo"
)

func main() {
	fmt.Println("Starting...")

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	viper.SetConfigName("salmon")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	port := viper.GetString("port")
	clean_db := viper.GetBool("clean_db")
	db_location := viper.GetString("db")

	if clean_db {
		fmt.Println("Initializing db")
		os.Remove(db_location)
	}

	db, err := sqlx.Open("sqlite3", db_location)

	if err != nil {
		panic(err)
	}

	if clean_db {
		db.MustExec(user_repo.Schema)
		db.MustExec(media_repo.Schema)
		db.MustExec(session_repo.Schema)
	}

	userRepo := user_repo.New(db)
	userService := user.CreateService(userRepo)

	mediaRepo := media_repo.New(db)
	mediaService := media.CreateService(mediaRepo)

	sessionRepo := session_repo.New(db)
	sessionService := session.CreateService(sessionRepo)

	if clean_db {
		userService.Signup("admin", "password", "email")
	}

	server := server.New(userService, mediaService, sessionService, port)

	if err = server.Run(); err != nil {
		return err
	}

	return nil
}
