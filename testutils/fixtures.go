package testutils

import "github.com/Alexande92/go-simple-library/internal/storage"

func GetTestBook() storage.Book {
	return storage.Book{
		Author:          "test",
		Title:           "test",
		PublicationDate: "2022-12",
		Publisher:       "test",
		Edition:         2,
		Location:        "test",
	}
}
