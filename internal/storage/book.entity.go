package storage

type Book struct {
	Id              int64  `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationDate string `json:"publication-date"`
	Publisher       string `json:"publisher,omitempty"`
	Edition         int    `json:"edition,omitempty"`
	Location        string `json:"location,omitempty"`
}

func (b *Book) getId() int64 {
	return b.Id
}

func (b *Book) setId(s Table) int64 {
	b.Id = int64(len(s) + 1)
	return b.Id
}
