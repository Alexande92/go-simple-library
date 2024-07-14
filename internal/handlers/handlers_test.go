package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Alexande92/go-simple-library/internal/entities"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"github.com/Alexande92/go-simple-library/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Path struct {
	mainRoute  string
	pathValues map[string]string
}

var apiUrl string = "/api/v1/books"

func sendTestRequest(t *testing.T, method string, path string, payload []byte, db *storage.Storage) (int, string, []byte) {
	t.Helper()

	body := bytes.NewReader(payload)

	router := http.NewServeMux()
	RegisterRoutes(router, db)
	srv := httptest.NewServer(router)

	req, err := http.NewRequest(method, srv.URL+path, body)
	require.NoError(t, err)

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	defer func() {
		_ = resp.Body.Close()
	}()

	require.NoError(t, err)

	contentType := resp.Header.Get("Content-Type")

	if resp.Body == http.NoBody {
		return resp.StatusCode, contentType, nil
	}

	respBody, err := io.ReadAll(resp.Body)

	require.NoError(t, err)

	return resp.StatusCode, contentType, respBody
}

func TestCheckHealth(t *testing.T) {
	code, contentType, body := sendTestRequest(t, http.MethodGet, "/api/v1/health", nil, nil)
	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal(contentType, "text/plain; charset=utf-8")
	assert.Equal("Healthy", string(body))
}

func TestGetAllBooksHandler_EmptyStorage(t *testing.T) {
	db := CreateTestStorage()
	code, contentType, body := sendTestRequest(t, http.MethodGet, apiUrl, nil, db)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, code)
	assert.Equal(contentType, "application/json")

	actualBody := strings.Trim(string(body), "\n")
	assert.Equal("[]", actualBody)
}

func TestGetAllBooksHandler_NotEmptyStorage(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())

	payload, err := json.Marshal(testutils.GetTestBook())
	require.NoError(t, err)

	code, contentType, body := sendTestRequest(t, http.MethodGet, apiUrl, payload, db)

	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal(contentType, "application/json")

	getEncodedBook := string(func() []byte {
		b := testutils.GetTestBook()
		b.Id = db.GetLastId()

		var buf bytes.Buffer

		err = json.NewEncoder(&buf).Encode([]entities.Book{b})
		assert.NoError(err)

		return buf.Bytes()
	}())

	assert.Equal(getEncodedBook, string(body))
}

func TestBookHandler_GetBookById_EmptyStorage(t *testing.T) {
	db := CreateTestStorage()
	code, contentType, body := sendTestRequest(t, http.MethodGet, apiUrl+"/1", nil, db)

	assert := assert.New(t)

	assert.Equal(http.StatusNotFound, code)
	assert.Equal(contentType, "application/json")

	assert.Equal("\"Book not found\"\n", string(body))
}

func TestBookHandler_GetBookById_NotEmptyStorage(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())
	code, contentType, body := sendTestRequest(t, http.MethodGet, apiUrl+"/1", nil, db)

	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal(contentType, "application/json")

	getEncodedBook := string(func() []byte {
		b := testutils.GetTestBook()
		b.Id = db.GetLastId()

		var buf bytes.Buffer

		err := json.NewEncoder(&buf).Encode(b)
		assert.NoError(err)

		return buf.Bytes()
	}())

	assert.Equal(getEncodedBook, string(body))
}

func TestBookHandler_GetBookByWrongId(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())
	code, contentType, body := sendTestRequest(t, http.MethodGet, apiUrl+"/test", nil, db)

	assert := assert.New(t)
	assert.Equal(contentType, "application/json")

	assert.Equal(http.StatusBadRequest, code)
	actualBody := strings.Trim(string(body), "\n")

	assert.Equal("\"Invalid book id\"", actualBody)
}

func TestBookHandler_DeleteBook(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())
	code, _, body := sendTestRequest(t, http.MethodDelete, apiUrl+"/1", nil, db)

	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal("", string(body))
}

func TestBookHandler_DeleteBookByWrongId(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())
	code, _, body := sendTestRequest(t, http.MethodDelete, apiUrl+"/test", nil, db)

	assert := assert.New(t)

	assert.Equal(http.StatusBadRequest, code)
	actualBody := strings.Trim(string(body), "\n")

	assert.Equal("\"Invalid book id\"", actualBody)
}

func TestBookHandler_SaveBookFailed_WrongJSON(t *testing.T) {
	db := CreateTestStorage()
	assert := assert.New(t)

	var buf bytes.Buffer

	err := json.NewEncoder(&buf).Encode("{{}")
	assert.NoError(err)

	code, contentType, body := sendTestRequest(t, http.MethodPost, apiUrl, buf.Bytes(), db)

	assert.Equal(http.StatusBadRequest, code)
	assert.Equal(contentType, "application/json")

	assert.Equal("\"Couldn't parse json\"\n", string(body))
}

func TestBookHandler_SaveBook(t *testing.T) {
	db := CreateTestStorage()
	book := testutils.GetTestBook()

	var buf bytes.Buffer
	assert := assert.New(t)

	err := json.NewEncoder(&buf).Encode(book)
	assert.NoError(err)
	code, contentType, body := sendTestRequest(t, http.MethodPost, apiUrl, buf.Bytes(), db)

	assert.Equal(http.StatusCreated, code)
	assert.Equal(contentType, "application/json")

	buf.Reset()
	book.Id = db.GetLastId()

	err = json.NewEncoder(&buf).Encode(book)

	assert.NoError(err)
	assert.Equal(buf.String(), string(body))
}

func TestBookHandler_SaveBook_ValidationError(t *testing.T) {
	db := CreateTestStorage()

	expected := ValidationErrors{
		Errors: []ErrorRes{
			{Field: "author", Reason: "field should be at least 1 characters long"},
			{Field: "publicationDate", Reason: "field should be equal to 7 chars"},
		},
	}

	book := testutils.GetTestBook()
	book.PublicationDate = "22-29"
	book.Author = ""

	var buf bytes.Buffer
	var encodedBuf bytes.Buffer

	json.NewEncoder(&buf).Encode(book)
	json.NewEncoder(&encodedBuf).Encode(expected)

	code, contentType, body := sendTestRequest(t, http.MethodPost, apiUrl, buf.Bytes(), db)

	assert := assert.New(t)

	assert.Equal(http.StatusBadRequest, code)
	assert.Equal(contentType, "application/json")

	assert.Equal(encodedBuf.String(), string(body))
}

func TestBookHandler_UpdateBook(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())

	book := testutils.GetTestBook()
	book.Id = db.GetLastId()
	book.Author = "A. Dyuma"
	var buf bytes.Buffer

	json.NewEncoder(&buf).Encode(book)

	code, contentType, body := sendTestRequest(t, http.MethodPut, apiUrl+"/1", buf.Bytes(), db)

	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal(contentType, "application/json")

	assert.Equal(buf.String(), string(body))
}

func TestBookHandler_UpdateBook_ValidationError(t *testing.T) {
	db := CreateTestStorage(testutils.GetTestBook())

	expected := ValidationErrors{
		Errors: []ErrorRes{
			{Field: "author", Reason: "field should be at least 1 characters long"},
			{Field: "publicationDate", Reason: "field should be equal to 7 chars"},
		},
	}

	book := testutils.GetTestBook()
	book.PublicationDate = "22-29"
	book.Author = ""

	var buf bytes.Buffer
	var encodedBuf bytes.Buffer

	json.NewEncoder(&buf).Encode(book)
	json.NewEncoder(&encodedBuf).Encode(expected)

	code, _, body := sendTestRequest(t, http.MethodPut, apiUrl+"/1", buf.Bytes(), db)

	assert := assert.New(t)

	assert.Equal(http.StatusBadRequest, code)
	assert.Equal(encodedBuf.String(), string(body))
}

func TestValidateBook(t *testing.T) {
	t.Run("Should be validated ok", func(t *testing.T) {
		book := testutils.GetTestBook()

		assert := assert.New(t)

		actual := ValidateBook(book)
		assert.Equal([]ErrorRes{}, actual)
	})

	t.Run("Should fail without required field", func(t *testing.T) {
		expected := errors.New("field should be at least 1 characters long")

		book := testutils.GetTestBook()
		book.PublicationDate = ""

		assert := assert.New(t)

		actual := validateLength(book.PublicationDate, 1, math.MaxInt)
		assert.ErrorAs(actual, &expected)
	})

	t.Run("Should fail when field not equal to 7 chars", func(t *testing.T) {
		expected := errors.New("field should be equal to 7 chars")

		book := testutils.GetTestBook()
		book.PublicationDate = "21-21"

		assert := assert.New(t)

		actual := isEqual(len(book.PublicationDate), 7)
		assert.ErrorAs(actual, &expected)
	})

	t.Run("Should fail when length more than 10 chars", func(t *testing.T) {
		expected := errors.New("field length should not exceed 10 characters")

		book := testutils.GetTestBook()
		book.Author = "Joan Joan Rouling"

		assert := assert.New(t)

		actual := validateLength(book.Author, 1, 10)
		assert.ErrorAs(actual, &expected)
	})

	t.Run("Should fail when length less than 10 chars", func(t *testing.T) {
		expected := errors.New("field length should not exceed 8 characters")

		book := testutils.GetTestBook()
		book.Author = "J.Rouling"

		assert := assert.New(t)

		actual := validateLength(book.Author, 1, 8)
		assert.ErrorAs(actual, &expected)
	})

	t.Run("Should be ok when length in proper range", func(t *testing.T) {

		book := testutils.GetTestBook()

		assert := assert.New(t)

		actual := validateLength(book.Author, 1, 255)
		assert.ErrorIs(nil, actual)
	})
}

func TestAddingValidationError(t *testing.T) {
	t.Run("Should return empty errors slice", func(t *testing.T) {
		actual := addValidationError([]ErrorRes{}, "test", nil)

		assert := assert.New(t)
		assert.Equal([]ErrorRes{}, actual)
	})

	t.Run("Should return error", func(t *testing.T) {
		err := errors.New("field should be at least 1 characters long")
		actual := addValidationError([]ErrorRes{}, "publicationDate", err)

		expected := []ErrorRes{
			{
				Field:  "publicationDate",
				Reason: err.Error(),
			},
		}

		assert := assert.New(t)
		assert.Equal(expected, actual)
	})
}

func CreateTestStorage(books ...entities.Book) *storage.Storage {
	db := storage.NewStorage()

	if len(books) > 0 {
		for _, book := range books {
			book = db.AddBook(book)

			db.Save(book)
		}
	}

	return db
}
