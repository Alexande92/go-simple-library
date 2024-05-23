package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_GetAll(t *testing.T) {
	storage := NewStorage()

	books := storage.GetAll()

	assert := assert.New(t)

	assert.Equal([]Book{}, books)
}

func TestStorage_Save(t *testing.T) {
	storage := NewStorage()
	book := getTestBook()

	storage.Save(book)

	t.Run("Last id should be equal to 1", func(t *testing.T) {
		assert := assert.New(t)
		assert.Equal(1, storage.lastId)
	})

	t.Run("Saved book should be equal to sent one", func(t *testing.T) {
		assert := assert.New(t)
		book.Id = 1
		assert.Equal(book, storage.books[storage.lastId])
	})
}

func TestStorage_Delete(t *testing.T) {
	storage := NewStorage()
	book := getTestBook()

	storage.Save(book)

	t.Run("Storage should be empty", func(t *testing.T) {
		assert := assert.New(t)
		storage.Delete(storage.lastId)

		assert.Equal([]Book{}, storage.GetAll())
	})
}

func TestStorage_GetById(t *testing.T) {
	storage := NewStorage()
	book := getTestBook()

	storage.Save(book)

	t.Run("Storage should return book by id", func(t *testing.T) {
		assert := assert.New(t)
		foundBook, _ := storage.GetById(storage.lastId)

		book.Id = 1
		assert.Equal(book, foundBook)
	})

	t.Run("Storage should return Err not found", func(t *testing.T) {
		assert := assert.New(t)
		_, err := storage.GetById(2)

		assert.Equal(ErrNotFound.Error(), err.Error())
	})

}

func getTestBook() Book {
	return Book{
		Author:          "test",
		Title:           "test",
		PublicationDate: "2022-12",
		Publisher:       "test",
		Edition:         2,
		Location:        "test",
	}
}
