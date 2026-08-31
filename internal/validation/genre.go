package validation

import (
	"errors"
	"movies-api/internal/models"
	"strings"
)

func ValidateGenre(genre models.Genre) error {
	if strings.TrimSpace(genre.Name) == "" {
		return errors.New("name is required")
	}

	if len(genre.Name) > 100 {
		return errors.New("name must be 100 characters or less")
	}

	return nil
}
