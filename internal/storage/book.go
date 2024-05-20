package storage

type Book struct {
	Id              int    `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationDate string `json:"publication-date"`
	Publisher       string `json:"publisher,omitempty"`
	Edition         int    `json:"edition,omitempty"`
	Location        string `json:"location,omitempty"`
}

func (b *Book) getId() int {
	return b.Id
}

func (b *Book) setId(s Storage) int {
	lastId := s.lastId + 1

	b.Id = lastId
	return b.Id
}
