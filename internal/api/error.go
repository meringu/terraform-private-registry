package api

import (
	"fmt"
	"net/http"
)

// ErrorResponse is the type returned on error
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// WriteError writes errors onto the response writer
func WriteError(w http.ResponseWriter, status int, errors ...error) error {
	errs := []string{}
	for _, err := range errors {
		errs = append(errs, err.Error())
	}
	return WriteJSON(w, status, ErrorResponse{
		Errors: errs,
	})
}

// WriteNotFound writes a not found error onto the response writer
func WriteNotFound(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotFound, fmt.Errorf("Not Found"))
}

// WriteNotImplemented writes a not implemented error onto the response writer
func WriteNotImplemented(w http.ResponseWriter) error {
	return WriteError(w, http.StatusNotImplemented, fmt.Errorf("Not Implemented"))
}
