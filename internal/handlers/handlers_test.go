package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetAllBooksHandler(t *testing.T) {
	db := storage.InitStorage()
	h := NewBaseHandler(db)

	//book := storage.Book{
	//	Author:          "test",
	//	Title:           "test",
	//	PublicationDate: "2022-12",
	//	Publisher:       "test",
	//	Edition:         2,
	//	Location:        "test",
	//}

	book := getTestBook()

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(book)

	req, _ := http.NewRequest("GET", "/api/v1/books", nil)
	saveReq, _ := http.NewRequest("POST", "/api/v1/book", &buf)

	res := httptest.NewRecorder()
	saveRes := httptest.NewRecorder()

	handler := http.HandlerFunc(h.GetBooks)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)
	assert.Equal(t, res.Body.String(), "[]\n")

	saveHandler := http.HandlerFunc(h.SaveBook)
	saveHandler.ServeHTTP(saveRes, saveReq)

	assert.Equal(t, 200, res.Code)

	book.Id = 1

	json.NewEncoder(&buf).Encode(book)
	assert.Equal(t, saveRes.Body.String(), buf.String())
}

func TestGetByIdBookHandler(t *testing.T) {
	db := storage.InitStorage()
	h := NewBaseHandler(db)

	req, _ := http.NewRequest("GET", "/api/v1/book/", nil)
	req.SetPathValue("id", "1")
	res := httptest.NewRecorder()

	handler := http.HandlerFunc(h.GetBookById)
	handler.ServeHTTP(res, req)

	assert.Equal(t, 404, res.Code)

	req, _ = http.NewRequest("GET", "/api/v1/book/", nil)
	req.SetPathValue("id", "tesd")

	res = httptest.NewRecorder()
	handler.ServeHTTP(res, req)

	assert.Equal(t, 400, res.Code)
	assert.Equal(t, res.Body.String(), "\"Could not fetch book, wrong id format\"\n")
	//assert.Equal(t, res.Body.String(), "[]\n")

	saveBookForTests(h)

	req, _ = http.NewRequest("GET", "/api/v1/book/1", nil)
	req.SetPathValue("id", "1")
	res = httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, 200, res.Code)

	var buf bytes.Buffer

	book := getTestBook()
	book.Id = 1

	json.NewEncoder(&buf).Encode(book)

	assert.Equal(t, res.Body.String(), buf.String())
}

func TestSaveBookHandler(t *testing.T) {
	db := storage.InitStorage()
	h := NewBaseHandler(db)

	book := getTestBook()

	res, _ := saveBookForTests(h)

	assert.Equal(t, 201, res.Code)

	var buf bytes.Buffer

	book.Id = 1

	json.NewEncoder(&buf).Encode(book)
	assert.Equal(t, res.Body.String(), buf.String())
}

func saveBookForTests(h *BaseHandler) (*httptest.ResponseRecorder, *http.Request) {
	book := getTestBook()

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(book)

	saveReq, _ := http.NewRequest("POST", "/api/v1/book", &buf)

	saveRes := httptest.NewRecorder()

	saveHandler := http.HandlerFunc(h.SaveBook)
	saveHandler.ServeHTTP(saveRes, saveReq)

	return saveRes, saveReq
}

func getTestBook() storage.Book {
	return storage.Book{
		Author:          "test",
		Title:           "test",
		PublicationDate: "2022-12",
		Publisher:       "test",
		Edition:         2,
		Location:        "test",
	}
}
