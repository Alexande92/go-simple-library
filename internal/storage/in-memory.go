package storage

// TODO make this storage not only book related, if needed

type Book struct {
	Id              int64  `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationDate string `json:"publication-date"`
	Publisher       string `json:"publisher,omitempty"`
	Edition         int    `json:"edition,omitempty"`
	Location        string `json:"location,omitempty"`
}

var BookStorage = map[int]Book{}

func (b *Book) Save() {
	b.Id = int64(len(BookStorage) + 1)
	BookStorage[len(BookStorage)+1] = *b
}

func GetAll() []Book {
	books := make([]Book, 0, len(BookStorage))

	for _, v := range BookStorage {
		books = append(books, v)
	}
	return books
}

func GetById(id int64) (Book, bool) {
	book, ok := BookStorage[int(id)]
	return book, ok
}

func Delete(id int64) {
	delete(BookStorage, int(id))
}

func (b *Book) Update() {
	BookStorage[int(b.Id)] = *b
}
