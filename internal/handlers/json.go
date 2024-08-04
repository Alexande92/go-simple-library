package handlers

type JsonErrorResponse struct {
	Message string `json:"message"`
}

type JsonValidationErrorResponse struct {
	Errors []ErrorRes `json:"errors,omitempty"`
}
