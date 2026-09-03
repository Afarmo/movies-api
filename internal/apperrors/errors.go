package apperrors

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidInput = errors.New("invalid input")
	ErrAssociated   = errors.New("associated")
	ErrDuplicateID  = errors.New("duplicate")
)

func WriteError(w http.ResponseWriter, err error) {

	switch {
	case errors.Is(err, ErrNotFound):
		http.Error(w, "Record Not Found", http.StatusNotFound)

	case errors.Is(err, ErrConflict):
		http.Error(w, "Conflict", http.StatusBadRequest)

	case errors.Is(err, ErrInvalidInput):
		http.Error(w, "Invalid Input", http.StatusBadRequest)

	case errors.Is(err, ErrAssociated):
		http.Error(w, err.Error(), http.StatusBadRequest)

	case errors.Is(err, ErrDuplicateID):
		http.Error(w, "Duplicate ID(s)", http.StatusBadRequest)

	default:
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}
