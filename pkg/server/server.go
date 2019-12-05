package server

import (
	"fmt"
	"github.com/steams/salmon-hub/pkg/media"
	"github.com/steams/salmon-hub/pkg/server/handlers"
	"github.com/steams/salmon-hub/pkg/user"
	"net/http"
)

type Server interface {
	Run() error
}

type server_imp struct {
	userService  user.Service
	mediaService media.Service
	port         string
}

func New(u user.Service, m media.Service, port string) Server {
	return server_imp{u, m, port}
}

func (s server_imp) Run() error {
	// user.Add("admin", "password", "email")
	// fmt.Println(user.Get("admin", "password"))

	http.HandleFunc("/media", handlers.Media_handler(s.mediaService))
	http.HandleFunc("/synch", handlers.Synch_handler(s.mediaService))
	http.HandleFunc("/api/login", handlers.Login_handler(s.userService))
	http.HandleFunc("/api/signup", handlers.Signup_handler(s.userService))
	// http.HandleFunc("/api/verify", verification_handler)
	http.HandleFunc("/", file_handler)
	return http.ListenAndServe(":"+s.port, nil)
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	enableCors(&w)

	fmt.Println(r.URL.Path)
	if r.URL.Path == "/elm.js" {
		http.ServeFile(w, r, "/home/steams/Development/audigo/salmon-web-client/elm.js")
		return
	}
	// path := "/home/steams/Development/audigo/salmon-web-client/" + r.URL.Path[1:]
	path := "/home/steams/Development/audigo/salmon-web-client/index.html"
	http.ServeFile(w, r, path)
}
