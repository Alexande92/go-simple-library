package storage

import (
	"github.com/Alexande92/go-simple-library/internal/entity"
	"github.com/Alexande92/go-simple-library/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStorage_GetAll_EmptyStorage(t *testing.T) {
	storage := NewStorage()
	books := storage.GetAll()

	assert.Equal(t, []entity.Book{}, books)
}

func TestStorage_GetAll_NotEmptyStorage(t *testing.T) {
	book := testutils.GetTestBook()

	storage := NewWithBooks(book)

	book.ID = storage.lastId
	books := storage.GetAll()
	assert.Equal(t, []entity.Book{
		book,
	}, books)
}

func TestStorage_Save(t *testing.T) {
	book := testutils.GetTestBook()
	storage := NewWithBooks(book)

	assert := assert.New(t)
	assert.Equal(1, storage.lastId)

	savedBook, err := storage.GetById(1)
	require.NoError(t, err)

	book.ID = storage.lastId
	assert.Equal(book, savedBook)
}

func TestStorage_Update(t *testing.T) {
	book := testutils.GetTestBook()
	storage := NewWithBooks(book)

	assert := assert.New(t)

	book, err := storage.GetById(storage.lastId)
	require.NoError(t, err)

	book.Author = "Duma Junior"
	err = storage.Update(book)
	require.NoError(t, err)

	updated, err := storage.GetById(book.ID)
	require.NoError(t, err)

	assert.Equal(book.Author, updated.Author)
}

func TestStorage_DeleteLastItem(t *testing.T) {
	storage := NewStorage()
	book := testutils.GetTestBook()

	storage.Save(book)

	err := storage.Delete(1)
	require.NoError(t, err)

	assert.Len(t, storage.GetAll(), 0)
}

func TestStorage_GetById(t *testing.T) {
	book := testutils.GetTestBook()
	storage := NewWithBooks(book)

	foundBook, err := storage.GetById(storage.lastId)
	book.ID = storage.lastId
	require.NoError(t, err)

	assert.Equal(t, book, foundBook)
}

func TestStorage_GetById_NotFound(t *testing.T) {
	storage := NewStorage()
	_, err := storage.GetById(2)
	assert.ErrorIs(t, ErrNotFound, err)
}
