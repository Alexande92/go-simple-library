package storage

import "errors"

// TODO make this storage not only book related, if needed

type Entity interface {
	getId() int
	setId(s Storage) int
}

var ErrNotFound = errors.New("entity not found")

type Storage struct {
	books  map[int]Book
	lastId int
}

func NewStorage() Storage {
	return Storage{}
}

func (s Storage) Save(e Book) {
	id := e.setId(s)
	s.books[id] = e
}

func (s Storage) GetAll() []Book {
	books := make([]Book, 0, len(s.books))

	for _, v := range s.books {
		books = append(books, v)
	}
	return books
}

func (s Storage) GetById(id int) (Book, error) {
	book, ok := s.books[id]

	if !ok {
		return Book{}, ErrNotFound
	}
	return book, nil
}

func (s Storage) Delete(id int) {
	delete(s.books, id)
}

func (s Storage) Update(e Book) {
	s.books[e.Id] = e
}
