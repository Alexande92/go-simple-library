package storage

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

func (b *Book) setId(s Storage) int {
	lastId := s.lastId + 1

	b.Id = lastId
	return b.Id
}

//func FromJson(data io.ReadCloser) (Book, error) {
//	var book Book
//
//	err := json.NewDecoder(data).Decode(&book)
//
//	if err != nil && err != io.EOF {
//		//w.WriteHeader(http.StatusInternalServerError)
//		//json.NewEncoder(w).Encode("Could not parse json")
//		return Book{}, errors.New("could not parse json")
//	}

//validatedErrs := handlers.ValidateBookData(book)

//if len(validatedErrs) != 0 {
//	return Book{},
//	//w.WriteHeader(http.StatusBadRequest)
//	err = json.NewEncoder(w).Encode(ValidationErrorRes{
//		Code:    400,
//		Message: "Validation failed",
//		Errors:  validatedErrs,
//	})
//
//	fmt.Println(err)
//	return
//}
//}
