package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alexande92/go-simple-library/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type Path struct {
	mainRoute  string
	pathValues map[string]string
}

func sendTestRequest(t *testing.T, method string, path Path, payload []byte, handler http.HandlerFunc) (int, string, []byte) {
	t.Helper()

	body := bytes.NewReader(payload)

	req, err := http.NewRequest(method, path.mainRoute, body)

	if len(path.pathValues) != 0 {
		for k, v := range path.pathValues {
			req.SetPathValue(k, v)
		}
	}

	require.NoError(t, err)
	respRecorder := httptest.NewRecorder()
	handler.ServeHTTP(respRecorder, req)

	require.NoError(t, err)

	if err != nil {
		fmt.Println(err.Error())
	}

	resp := respRecorder.Result()

	defer func() {
		_ = resp.Body.Close()
	}()

	contentType := resp.Header.Get("Content-Type")

	if resp.Body == http.NoBody {
		return resp.StatusCode, contentType, nil
	}

	respBody, err := io.ReadAll(resp.Body)

	require.NoError(t, err)

	return resp.StatusCode, contentType, respBody
}

func TestCheckHealth(t *testing.T) {
	code, _, body := sendTestRequest(t, http.MethodGet,
		Path{
			mainRoute: "/api/v1/health",
		}, nil, CheckHealth)

	assert := assert.New(t)

	assert.Equal(http.StatusOK, code)
	assert.Equal("Healthy", string(body))
}

func TestGetAllBooksHandler(t *testing.T) {
	db := storage.NewStorage()
	h := NewBookHandler(db)

	apiPath := Path{
		mainRoute: "/api/v1/books",
	}

	t.Run("Should get no books", func(t *testing.T) {
		code, _, body := sendTestRequest(t, http.MethodGet, apiPath, nil, h.GetBooks)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, code)
		actualBody := strings.Trim(string(body), "\n")
		assert.Equal("[]", actualBody)

	})

	t.Run("Should get all books", func(t *testing.T) {
		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodGet, apiPath, nil, h.GetBooks)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, code)

		getEncodedBook := string(func() []byte {
			b := getTestBook()
			var buf bytes.Buffer

			b.Id = 1

			json.NewEncoder(&buf).Encode([]storage.Book{b})
			return buf.Bytes()
		}())

		assert.Equal(getEncodedBook, string(body))
	})
}

func TestBookHandler_GetBookById(t *testing.T) {
	db := storage.NewStorage()
	h := NewBookHandler(db)

	apiPath := Path{
		mainRoute:  "/api/v1/books",
		pathValues: map[string]string{"id": "1"},
	}

	t.Run("Should find no books", func(t *testing.T) {
		code, _, body := sendTestRequest(t, http.MethodGet, apiPath, nil, h.GetBookById)

		assert := assert.New(t)

		assert.Equal(http.StatusNotFound, code)
		assert.Equal("\"book not found\"\n", string(body))

	})

	t.Run("Should get first book", func(t *testing.T) {
		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodGet, apiPath, nil, h.GetBookById)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, code)
		getEncodedBook := string(func() []byte {
			b := getTestBook()
			var buf bytes.Buffer

			b.Id = 1

			json.NewEncoder(&buf).Encode(b)
			return buf.Bytes()
		}())

		assert.Equal(getEncodedBook, string(body))
	})

	apiPath.pathValues["id"] = "test"

	t.Run("Should get error with wrong id", func(t *testing.T) {
		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodGet, apiPath, nil, h.GetBookById)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, code)
		actualBody := strings.Trim(string(body), "\n")

		assert.Equal("\"Invalid book id\"", actualBody)
	})
}

func TestBookHandler_DeleteBook(t *testing.T) {
	db := storage.NewStorage()
	h := NewBookHandler(db)

	apiPath := Path{
		mainRoute:  "/api/v1/books",
		pathValues: map[string]string{"id": "1"},
	}

	t.Run("Should delete book", func(t *testing.T) {
		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodDelete, apiPath, nil, h.DeleteBook)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, code)
		assert.Equal("", string(body))
	})

	apiPath.pathValues["id"] = "test"

	t.Run("Should get error with wrong id", func(t *testing.T) {
		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodDelete, apiPath, nil, h.DeleteBook)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, code)
		actualBody := strings.Trim(string(body), "\n")

		assert.Equal("\"Invalid book id\"", actualBody)
	})
}

func TestBookHandler_SaveBook(t *testing.T) {
	db := storage.NewStorage()
	h := NewBookHandler(db)

	apiPath := Path{
		mainRoute: "/api/v1/books",
	}

	t.Run("Should fail when sent wrong json", func(t *testing.T) {
		var buf bytes.Buffer

		json.NewEncoder(&buf).Encode("{{}")
		code, _, body := sendTestRequest(t, http.MethodPost, apiPath, buf.Bytes(), h.SaveBook)

		assert := assert.New(t)

		assert.Equal(http.StatusInternalServerError, code)
		assert.Equal("\"Could not parse json\"\n", string(body))
	})

	t.Run("Should save book", func(t *testing.T) {
		book := getTestBook()
		var buf bytes.Buffer

		json.NewEncoder(&buf).Encode(book)
		code, _, body := sendTestRequest(t, http.MethodPost, apiPath, buf.Bytes(), h.SaveBook)

		assert := assert.New(t)

		assert.Equal(http.StatusCreated, code)
		assert.Equal(buf.String(), string(body))
	})

	t.Run("Should fail validation", func(t *testing.T) {
		expected := ValidationErrors{
			Errors: []ErrorRes{
				{Field: "author", Reason: "missing required field"},
				{Field: "publicationDate", Reason: "field should be equal to 7 chars"},
			},
		}

		book := getTestBook()
		book.PublicationDate = "22-29"
		book.Author = ""

		var buf bytes.Buffer
		var encodedBuf bytes.Buffer

		json.NewEncoder(&buf).Encode(book)
		json.NewEncoder(&encodedBuf).Encode(expected)

		code, _, body := sendTestRequest(t, http.MethodPost, apiPath, buf.Bytes(), h.SaveBook)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, code)
		assert.Equal(encodedBuf.String(), string(body))
	})
}

func TestBookHandler_UpdateBook(t *testing.T) {
	db := storage.NewStorage()
	h := NewBookHandler(db)

	apiPath := Path{
		mainRoute:  "/api/v1/books",
		pathValues: map[string]string{"id": "1"},
	}

	t.Run("Should update book", func(t *testing.T) {
		book := getTestBook()
		book.Id = 1
		book.Author = "A. Dyuma"
		var buf bytes.Buffer

		json.NewEncoder(&buf).Encode(book)

		h.db.Save(getTestBook())
		code, _, body := sendTestRequest(t, http.MethodPut, apiPath, buf.Bytes(), h.UpdateBook)

		assert := assert.New(t)

		assert.Equal(http.StatusOK, code)
		assert.Equal(buf.String(), string(body))
	})

	t.Run("Should fail validation", func(t *testing.T) {
		expected := ValidationErrors{
			Errors: []ErrorRes{
				{Field: "author", Reason: "missing required field"},
				{Field: "publicationDate", Reason: "field should be equal to 7 chars"},
			},
		}

		book := getTestBook()
		book.PublicationDate = "22-29"
		book.Author = ""

		var buf bytes.Buffer
		var encodedBuf bytes.Buffer

		json.NewEncoder(&buf).Encode(book)
		json.NewEncoder(&encodedBuf).Encode(expected)

		code, _, body := sendTestRequest(t, http.MethodPost, apiPath, buf.Bytes(), h.UpdateBook)

		assert := assert.New(t)

		assert.Equal(http.StatusBadRequest, code)
		assert.Equal(encodedBuf.String(), string(body))
	})
}

func TestValidateBook(t *testing.T) {
	t.Run("Should be validated ok", func(t *testing.T) {
		book := getTestBook()

		assert := assert.New(t)

		actual := ValidateBook(book)
		assert.Equal([]ErrorRes{}, actual)
	})

	t.Run("Should fail without required field", func(t *testing.T) {
		expected := errors.New("missing required field")

		book := getTestBook()
		book.PublicationDate = ""

		assert := assert.New(t)

		actual := validateEmptiness(book.PublicationDate)
		assert.Equal(expected, actual)
	})

	t.Run("Should fail when field not equal to 7 chars", func(t *testing.T) {
		expected := errors.New("field should be equal to 7 chars")

		book := getTestBook()
		book.PublicationDate = "21-21"

		assert := assert.New(t)

		actual := isEqual(len(book.PublicationDate), 7)
		assert.Equal(expected, actual)
	})

	t.Run("Should fail when length more than 10 chars", func(t *testing.T) {
		expected := errors.New("field should contain less or equal than 10 chars")

		book := getTestBook()
		book.Author = "Joan Joan Rouling"

		assert := assert.New(t)

		actual := validateLength(len(book.Author), -1, 10)
		assert.Equal(expected, actual)
	})

	t.Run("Should fail when length less than 10 chars", func(t *testing.T) {
		expected := errors.New("field should contain more or equal than 10 chars")

		book := getTestBook()
		book.Author = "J.Rouling"

		assert := assert.New(t)

		actual := validateLength(len(book.Author), 10, -1)
		assert.Equal(expected, actual)
	})

	t.Run("Should be ok when length in proper range", func(t *testing.T) {

		book := getTestBook()

		assert := assert.New(t)

		actual := validateLength(len(book.Author), 1, 255)
		assert.Equal(nil, actual)
	})
}

func TestAddingValidationError(t *testing.T) {
	t.Run("Should return empty errors slice", func(t *testing.T) {
		actual := addValidationError([]ErrorRes{}, "test", nil)

		assert := assert.New(t)
		assert.Equal([]ErrorRes{}, actual)
	})

	t.Run("Should return error", func(t *testing.T) {
		err := errors.New("missing required field")
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
