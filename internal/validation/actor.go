package validation

import (
	"errors"
	"movies-api/internal/models"
	"strings"
	"time"
)

func ValidateActor(actor models.Actor) error {
	if strings.TrimSpace(actor.Name) == "" {
		return errors.New("name is required")
	}

	if len(actor.Name) > 200 {
		return errors.New("name must be 200 characters or less")
	}

	if strings.TrimSpace(actor.BirthDate) == "" {
		return errors.New("birthdate is required")
	}

	_, err := time.Parse("2006-01-02", actor.BirthDate)
	if err != nil {
		return errors.New("birthdate must be a valid date in YYYY-MM-DD format")
	}

	if actor.BirthDate > time.Now().Format("2006-01-02") {
		return errors.New("birthdate cannot be in the future")
	}

	return nil
}
