package storage

import (
	"errors"
	"github.com/Alexande92/go-simple-library/internal/entity"
)

var ErrNotFound = errors.New("book not found")

type Storage struct {
	books  map[int]entity.Book
	lastId int
}

func NewStorage() *Storage {
	return &Storage{
		books: make(map[int]entity.Book),
	}
}

func NewWithBooks(books ...entity.Book) *Storage {
	db := NewStorage()

	if len(books) > 0 {
		for _, book := range books {
			book = db.Save(book)
		}
	}

	return db
}

func (s *Storage) Save(b entity.Book) entity.Book {
	s.lastId++
	b.ID = s.lastId

	s.books[b.ID] = b
	return b
}

func (s *Storage) GetAll() []entity.Book {
	books := make([]entity.Book, 0, len(s.books))

	for _, v := range s.books {
		books = append(books, v)
	}
	return books
}

func (s *Storage) GetById(id int) (entity.Book, error) {
	book, ok := s.books[id]

	if !ok {
		return entity.Book{}, ErrNotFound
	}
	return book, nil
}

func (s *Storage) Delete(id int) error {
	if _, ok := s.books[id]; !ok {
		return ErrNotFound
	}

	delete(s.books, id)
	return nil
}

func (s *Storage) Update(b entity.Book) error {
	if _, ok := s.books[b.ID]; !ok {
		return ErrNotFound
	}

	s.books[b.ID] = b
	return nil
}
