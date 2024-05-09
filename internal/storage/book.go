package storage

import (
	"math/rand"
	"time"
)

type Book struct {
	Id              int    `json:"id"`
	Author          string `json:"author"`
	Title           string `json:"title"`
	PublicationDate string `json:"publication-date"`
	Publisher       string `json:"publisher,omitempty"`
	Edition         int    `json:"edition,omitempty"`
	Location        string `json:"location,omitempty"`
}

func (b *Book) getId() int {
	return b.Id
}

func (b *Book) setId(s Table) int {
	seed := rand.NewSource(time.Now().UnixNano())
	randomizer := rand.New(seed)

	b.Id = randomizer.Intn(1000)
	return b.Id
}
