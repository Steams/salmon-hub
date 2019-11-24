package media

import (
	"github.com/steams/salmon-hub/pkg/storage"
)

type Media struct {
	Title    string
	Artist   string
	Album    string
	Duration int
	Playlist string
	Hash     string
}

func Add(id string, m Media) {
	storage.MediaStore().Add(to_db(id, m))
}

func List(id string) []Media {
	return from_db_list(storage.MediaStore().List(id))
}

func ListHashes(id string) []string {
	return storage.MediaStore().ListHashes(id)
}

func Delete(id, hash string) {
	storage.MediaStore().Delete(id, hash)
}

func from_db(m storage.Media) Media {
	return Media{
		Title:    m.Title,
		Artist:   m.Artist,
		Album:    m.Album,
		Duration: m.Duration,
		Playlist: m.Playlist,
		Hash:     m.Hash,
	}
}

func to_db(id string, m Media) storage.Media {
	return storage.Media{
		Title:    m.Title,
		Artist:   m.Artist,
		Album:    m.Album,
		Duration: m.Duration,
		Hash:     m.Hash,
		Playlist: m.Playlist,
		UserID:   id,
	}
}

func from_db_list(m []storage.Media) []Media {
	var media []Media

	for _, song := range m {
		media = append(media, from_db(song))
	}

	return media
}
