package server

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/storage"
	"log"
	"net/http"
	"os"
)

func media_handler(store storage.MediaStorage) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case "GET":

			media_list := store.List()

			js, err := json.Marshal(media_list)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)

		case "POST":
			var m storage.Media

			fmt.Println("POST REQUEST")
			if r.Body == nil {
				http.Error(w, "Please send a request body", 400)
				return
			}
			err := json.NewDecoder(r.Body).Decode(&m)
			fmt.Fprintln(os.Stdout, "BODY: \n%+v", m)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			store.Add(m)
		}
	}

}

func file_handler(w http.ResponseWriter, r *http.Request) {
	path := "/home/steams/Development/audigo/salmon-web-client/" + r.URL.Path[1:]
	http.ServeFile(w, r, path)
}

func Run() {

	storage := storage.Open()

	http.HandleFunc("/", file_handler)
	http.HandleFunc("/media", media_handler(storage))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
