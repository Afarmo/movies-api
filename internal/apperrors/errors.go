package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidInput = errors.New("invalid input")
)

func WriteError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, "Record Not Found", http.StatusNotFound)

	case errors.Is(err, ErrConflict):
		http.Error(w, "Conflict", http.StatusConflict)

	case errors.Is(err, ErrInvalidInput):
		http.Error(w, "Invalid Input", http.StatusBadRequest)

	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
