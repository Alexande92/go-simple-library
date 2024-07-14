package storage

import (
	"errors"
	"github.com/Alexande92/go-simple-library/internal/entities"
)

var ErrNotFound = errors.New("book not found")

type Storage struct {
	books  map[int]entities.Book
	lastId int
}

func NewStorage() *Storage {
	return &Storage{
		books: make(map[int]entities.Book),
	}
}

func (s *Storage) GetLastId() int {
	return s.lastId
}

func (s *Storage) AddBook(book entities.Book) entities.Book {
	s.lastId++
	book.Id = s.lastId
	return book
}

func (s *Storage) Save(b entities.Book) entities.Book {
	s.books[b.Id] = b
	return b
}

func (s *Storage) GetAll() []entities.Book {
	books := make([]entities.Book, 0, len(s.books))

	for _, v := range s.books {
		books = append(books, v)
	}
	return books
}

func (s *Storage) GetById(id int) (entities.Book, error) {
	book, ok := s.books[id]

	if !ok {
		return entities.Book{}, ErrNotFound
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

func (s *Storage) Update(b entities.Book) error {
	_, ok := s.books[b.Id]

	if !ok {
		return ErrNotFound
	}

	s.books[b.Id] = b
	return nil
}
