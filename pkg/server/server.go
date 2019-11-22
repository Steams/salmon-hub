package server

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/domain/media"
	"github.com/steams/salmon-hub/pkg/domain/user"
	"log"
	"net/http"
	"os"
)

func synch_handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case "GET":

		id := r.Header.Get("Authorization")

		if id == "" {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}
		hashes := media.ListHashes(id)

		js, err := json.Marshal(hashes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	case "POST":
		var hashes []string

		id := r.Header.Get("Authorization")
		if id == "" {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		fmt.Println("POST DELETE REQUEST")
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&hashes)
		fmt.Fprintln(os.Stdout, "HASHES: %+v", hashes)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, hash := range hashes {
			media.Delete(id, hash)
		}
	}

}

func media_handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

		id := r.Header.Get("Authorization")

		if id == "" {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		media_list := media.List(id)

		js, err := json.Marshal(media_list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
		return

	case "POST":

		id := r.Header.Get("Authorization")

		if id == "" {
			http.Error(w, "userid not present", http.StatusBadRequest)
			return
		}

		// NOTE Perhaps its idiomatic in go to call this media.Model instead ?
		var songs []media.Media

		fmt.Println("POST REQUEST")
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&songs)
		fmt.Fprintln(os.Stdout, "BODY: \n%+v", songs)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		for _, song := range songs {
			fmt.Println(song)
			media.Add(id, song)
		}
		w.WriteHeader(http.StatusOK)
	}
}

func file_handler(w http.ResponseWriter, r *http.Request) {
	path := "/home/steams/Development/audigo/salmon-web-client/" + r.URL.Path[1:]
	http.ServeFile(w, r, path)
}

func Run() {

	user.Add("admin", "password")
	fmt.Println(user.Get("admin", "password"))

	http.HandleFunc("/", file_handler)
	http.HandleFunc("/media", media_handler)
	http.HandleFunc("/synch", synch_handler)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
