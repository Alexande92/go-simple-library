package testutils

import (
	"github.com/Alexande92/go-simple-library/internal/entity"
)

func GetTestBook() entity.Book {
	return entity.Book{
		Author:          "test",
		Title:           "test",
		PublicationDate: "2022-12",
		Publisher:       "test",
		Edition:         2,
		Location:        "test",
	}
}
