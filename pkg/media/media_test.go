package media

import (
	"github.com/steams/salmon-hub/pkg/rand"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAdd(t *testing.T) {
	// Test that Add method on media service passes the apropriate values through to the repository

	repo := new_test_repo()
	media_service := CreateService(repo)

	input := Media{"title", "artist", "album", 1, "playlist", "hash"}
	id := "id"

	media_service.Add(id, input)

	assert.Equal(t, input, repo.added[id], "Media Service Should pass correct struct and id to repository")
}

func TestDelete(t *testing.T) {
	// Test that Add method on media service passes the apropriate values through to the repository

	song := Media{"title", "artist", "album", 1, "playlist", "hash"}
	id := "id"

	repo := new_test_repo()
	repo.Add(id, song)

	media_service := CreateService(repo)
	media_service.Delete(id, song.Hash)

	assert.Empty(t, repo.added)
}

func TestList(t *testing.T) {
	var songs []Media
	id := rand.String(10)

	for i := 0; i <= rand.Int(10); i++ {
		rand_song := Media{rand.String(10), rand.String(10), rand.String(10), rand.Int(100), rand.String(10), rand.String(10)}
		songs = append(songs, rand_song)
	}

	repo := new_test_repo()
	repo.list[id] = songs
	media_service := CreateService(repo)

	result := media_service.List(id)

	assert.Equal(t, songs, result, "Media Service Should return appropriate list of media contained in repository")

	// Test with empty repo
	repo = new_test_repo()
	media_service = CreateService(repo)
	result = media_service.List(id)

	// TODO should probably return error if id doesnt exist ?
	assert.Equal(t, []Media{}, result, "Media Service Should return empty list if repo contains no media")
}

func TestListHashes(t *testing.T) {
	var songs []Media
	var hashes []string
	id := rand.String(10)

	for i := 0; i <= rand.Int(10); i++ {
		hash := rand.String(10)
		rand_song := Media{Hash: hash}
		hashes = append(hashes, hash)
		songs = append(songs, rand_song)
	}

	repo := new_test_repo()
	repo.list[id] = songs
	media_service := CreateService(repo)

	result := media_service.ListHashes(id)

	assert.Equal(t, hashes, result, "Media Service Should return appropriate list of hashes contained in repository")

	// Test with empty repo
	repo = new_test_repo()
	media_service = CreateService(repo)
	result = media_service.ListHashes(id)

	assert.Equal(t, []string{}, result, "Media Service Should return empty list if repo contains no media")
}

func new_test_repo() test_repo {
	return test_repo{make(map[string]Media), make(map[string][]Media)}
}

type test_repo struct {
	added map[string]Media
	list  map[string][]Media
}

func (t test_repo) Add(id string, m Media) {
	t.added[id] = m
}

func (t test_repo) List(id string) []Media {
	return t.list[id]
}

func (t test_repo) ListHashes(id string) []string {
	var hashes []string
	for _, song := range t.list[id] {
		hashes = append(hashes, song.Hash)
	}
	return hashes
}

func (t test_repo) Delete(id, hash string) {
	if hash == t.added[id].Hash {
		t.added = nil
	}
}
