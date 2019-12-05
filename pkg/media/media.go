package media

type Media struct {
	Title    string
	Artist   string
	Album    string
	Duration int
	Playlist string
	Hash     string
}

type Repository interface {
	Add(id string, m Media)
	List(user_id string) []Media
	ListHashes(user_id string) []string
	Delete(user_id, hash string)
}

type Service interface {
	Add(user_id string, m Media)
	List(id string) []Media
	ListHashes(id string) []string
	Delete(id, hash string)
}

func CreateService(r Repository) Service {
	return service_imp{r}
}

type service_imp struct {
	repo Repository
}

func (s service_imp) Add(id string, m Media) {
	s.repo.Add(id, m)
}

func (s service_imp) List(id string) []Media {
	return s.repo.List(id)
}

func (s service_imp) ListHashes(id string) []string {
	return s.repo.ListHashes(id)
}

func (s service_imp) Delete(id, hash string) {
	s.repo.Delete(id, hash)
}
