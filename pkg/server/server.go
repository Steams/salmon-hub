package server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
	"github.com/steams/salmon-hub/pkg/media"
	"github.com/steams/salmon-hub/pkg/server/handlers"
	"github.com/steams/salmon-hub/pkg/session"
	"github.com/steams/salmon-hub/pkg/user"
)

type Server interface {
	Run() error
}

type server_imp struct {
	userService    user.Service
	mediaService   media.Service
	sessionService session.Service
	port           string
}

func New(u user.Service, m media.Service, s session.Service, port string) Server {
	return server_imp{u, m, s, port}
}

func (s server_imp) Run() error {
	// user.Add("admin", "password", "email")
	// fmt.Println(user.Resolve("admin", "password"))

	http.HandleFunc("/media", handlers.Media_handler(s.mediaService, s.sessionService))
	http.HandleFunc("/signup", handlers.Signup_handler(s.userService))
	http.HandleFunc("/api/login", handlers.Login_handler(s.userService, s.sessionService))
	http.HandleFunc("/csrf", handlers.Csrf_handler(s.userService, s.sessionService))
	// Routes for media server, these routes dont use cookie based auth
	http.HandleFunc("/api/synch", handlers.Synch_handler(s.mediaService))
	http.HandleFunc("/api/media", handlers.API_Media_handler(s.mediaService))
	http.HandleFunc("/api/register", handlers.Register_handler(s.userService, s.sessionService))
	// http.HandleFunc("/api/verify", verification_handler)
	http.HandleFunc("/web_assets/", assets_handler)
	http.HandleFunc("/", file_handler)
	return http.ListenAndServe(":"+s.port, nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println("fie handler url path")
	fmt.Println(r.URL.Path)
	if r.URL.Path == "/elm.js" {

		if viper.GetBool("dev") {
			http.ServeFile(w, r, "/home/steams/Development/audigo/salmon-web-client/elm.js")
			return
		}

		dir, _ := os.Getwd()
		path := dir + "/web_assets/elm.js"
		http.ServeFile(w, r, path)
		return
	}

	if viper.GetBool("dev") {

		fmt.Println("DEV MODE")
		path := "/home/steams/Development/audigo/salmon-web-client/index.html"
		http.ServeFile(w, r, path)
		return
	}

	fmt.Println("PROD MODE")

	dir, _ := os.Getwd()
	path := dir + "/web_assets/index.html"
	http.ServeFile(w, r, path)
}

func assets_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	dir, _ := os.Getwd()

	path := dir + r.URL.Path

	http.ServeFile(w, r, path)
}
