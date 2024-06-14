package storage

import (
	"github.com/Alexande92/go-simple-library/testutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStorage_GetAll_EmptyStorage(t *testing.T) {
	storage := NewStorage()
	books := storage.GetAll()

	assert := assert.New(t)

	//t.Run("Storage should be empty", func(t *testing.T) {
	assert.Equal([]Book{}, books)
	//})

	//t.Run("Storage should contain a test book", func(t *testing.T) {
	//	book := getTestBook()
	//	book = storage.Save(book)
	//
	//	books = storage.GetAll()
	//	assert.Equal([]Book{
	//		book,
	//	}, books)
	//})

}

func TestStorage_GetAll_NotEmptyStorage(t *testing.T) {
	book := testutils.GetTestBook()

	storage := testutils.CreateTestStorage(book)

	//book := getTestBook()
	//book = storage.Save(book)

	books := storage.GetAll()
	assert.Equal(t, []Book{
		book,
	}, books)
}

func TestStorage_Save(t *testing.T) {
	book := testutils.GetTestBook()
	storage := testutils.CreateTestStorage(book)
	//book := getTestBook()
	assert := assert.New(t)
	//book = storage.AddBook(book)
	//storage.Save(book)
	assert.Equal(1, storage.lastId)

	savedBook, err := storage.GetById(storage.GetLastId())
	assert.NoError(err)

	assert.Equal(book, savedBook)
}

//t.Run("Last id should be equal to 1", func(t *testing.T) {
//	assert.Equal(t, 1, storage.lastId)
//})
//
//t.Run("Saved book should have id equal to 2", func(t *testing.T) {
//	book = storage.Save(book)
//	assert.Equal(t, book.Id, 2)
//})

func TestStorage_Delete(t *testing.T) {
	storage := NewStorage()
	book := testutils.GetTestBook()

	storage.Save(book)

	//t.Run("Storage should be empty", func(t *testing.T) {
	storage.Delete(storage.GetLastId())

	assert.Equal(t, []Book{}, storage.GetAll())
	//})
}

func TestStorage_GetById(t *testing.T) {
	book := testutils.GetTestBook()
	storage := testutils.CreateTestStorage(book)

	//
	//book = storage.Save(book)

	t.Run("Storage should return book by id", func(t *testing.T) {
		foundBook, err := storage.GetById(storage.lastId)
		assert.NoError(t, err)
		assert.Equal(t, book, foundBook)
	})

	t.Run("Storage should return Err not found", func(t *testing.T) {
		_, err := storage.GetById(2)
		assert.ErrorIs(t, ErrNotFound, err)
	})

}

//func getTestBook() Book {
//	return Book{
//		Author:          "test",
//		Title:           "test",
//		PublicationDate: "2022-12",
//		Publisher:       "test",
//		Edition:         2,
//		Location:        "test",
//	}
//}
