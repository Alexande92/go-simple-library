package storage

import (
	"github.com/Alexande92/go-simple-library/internal/entities"
	"github.com/Alexande92/go-simple-library/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_GetAll_EmptyStorage(t *testing.T) {
	storage := NewStorage()
	books := storage.GetAll()

	assert := assert.New(t)
	assert.Equal([]entities.Book{}, books)
}

func TestStorage_GetAll_NotEmptyStorage(t *testing.T) {
	book := testutils.GetTestBook()

	storage := CreateTestStorage(book)

	book.Id = storage.lastId
	books := storage.GetAll()
	assert.Equal(t, []entities.Book{
		book,
	}, books)
}

func TestStorage_Save(t *testing.T) {
	book := testutils.GetTestBook()
	storage := CreateTestStorage(book)

	assert := assert.New(t)
	assert.Equal(1, storage.lastId)

	savedBook, err := storage.GetById(storage.GetLastId())
	assert.NoError(err)

	book.Id = storage.lastId
	assert.Equal(book, savedBook)
}

func TestStorage_Update(t *testing.T) {
	book := testutils.GetTestBook()
	storage := CreateTestStorage(book)

	assert := assert.New(t)

	book, err := storage.GetById(storage.lastId)
	assert.NoError(err)

	book.Author = "Duma Junior"
	err = storage.Update(book)
	assert.NoError(err)

	updated, err := storage.GetById(book.Id)
	assert.NoError(err)

	assert.Equal(book.Author, updated.Author)
}

func TestStorage_DeleteLastItem(t *testing.T) {
	storage := NewStorage()
	book := testutils.GetTestBook()

	storage.Save(book)

	err := storage.Delete(storage.GetLastId())
	assert.NoError(t, err)

	assert.Len(t, storage.GetAll(), 0)
}

func TestStorage_GetById(t *testing.T) {
	book := testutils.GetTestBook()
	storage := CreateTestStorage(book)

	t.Run("Storage should return book by id", func(t *testing.T) {
		foundBook, err := storage.GetById(storage.lastId)
		book.Id = storage.lastId
		assert.NoError(t, err)
		assert.Equal(t, book, foundBook)
	})

	t.Run("Storage should return Err not found", func(t *testing.T) {
		_, err := storage.GetById(2)
		assert.ErrorIs(t, ErrNotFound, err)
	})

}

func CreateTestStorage(books ...entities.Book) *Storage {
	db := NewStorage()

	if len(books) > 0 {
		for _, book := range books {
			book = db.AddBook(book)

			db.Save(book)
		}
	}

	return db
}
