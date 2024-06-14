package storage

import (
	"errors"
)

// TODO make this storage not only book related, if needed

type Entity interface {
	getId() int
	setId(s Storage) int
}

var ErrNotFound = errors.New("book not found")

type Storage struct {
	books  map[int]Book
	lastId int
}

func NewStorage() *Storage {
	return &Storage{
		books: make(map[int]Book),
	}
}

func (s *Storage) GetLastId() int {
	return s.lastId
}

func (s *Storage) AddBook(book Book) Book {
	s.lastId++
	book.Id = s.lastId
	return book
}

func (s *Storage) Save(b Book) Book {
	//s.lastId++
	//b.Id = s.lastId
	//id := b.setId(*s)
	s.books[b.Id] = b
	//s.lastId = id
	return b
}

func (s *Storage) GetAll() []Book {
	books := make([]Book, 0, len(s.books))

	for _, v := range s.books {
		books = append(books, v)
	}
	return books
}

func (s *Storage) GetById(id int) (Book, error) {
	book, ok := s.books[id]

	if !ok {
		return Book{}, ErrNotFound
	}
	return book, nil
}

func (s *Storage) Delete(id int) error {
	_, ok := s.books[id]

	if !ok {
		return ErrNotFound
	}

	delete(s.books, id)
	return nil
}

func (s *Storage) Update(b Book) error {
	_, ok := s.books[b.Id]

	if !ok {
		return ErrNotFound
	}

	s.books[b.Id] = b
	return nil
}
