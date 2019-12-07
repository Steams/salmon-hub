package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/steams/salmon-hub/pkg/media"
	"github.com/steams/salmon-hub/pkg/session"
	"net/http"
	"os"
)

func Synch_handler(service media.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "GET":
			id := r.Header.Get("Authorization")

			if id == "" {
				http.Error(w, "userid not present", http.StatusBadRequest)
				return
			}

			hashes := service.ListHashes(id)

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
				service.Delete(id, hash)
			}
		}

	}
}

func Media_handler(media_service media.Service, session_service session.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "GET":

			token, err := r.Cookie("session_token")

			if err != nil {
				http.Error(w, "userid not present", http.StatusBadRequest)
				return
			}

			id := session_service.Get(token.Value)

			media_list := media_service.List(id)

			if media_list == nil {
				js, _ := json.Marshal([]media.Media{})
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				return
			}

			js, err := json.Marshal(media_list)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
			return

		}
	}
}

func API_Media_handler(service media.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// defer enableCors(&w)

		switch r.Method {
		case "GET":

			id := r.Header.Get("Authorization")

			if id == "" {
				http.Error(w, "userid not present", http.StatusBadRequest)
				return
			}

			// media_list := service.List(id.Value)
			media_list := service.List(id)

			if media_list == nil {
				js, _ := json.Marshal([]media.Media{})
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
				return
			}

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
				// service.Add(id.Value, song)
				service.Add(id, song)
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}
