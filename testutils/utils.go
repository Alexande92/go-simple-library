package testutils

import "github.com/Alexande92/go-simple-library/internal/storage"

func CreateTestStorage(books ...storage.Book) *storage.Storage {
	db := storage.NewStorage()

	if len(books) > 0 {
		for _, book := range books {
			book = db.AddBook(book)

			db.Save(book)
		}
	}

	return db
}
