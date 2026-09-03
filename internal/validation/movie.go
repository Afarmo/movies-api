package validation

import (
	"errors"
	"fmt"
	"movies-api/internal/models"
	"strings"
	"time"
)

func ValidateMovie(movie models.Movie) error {
	if strings.TrimSpace(movie.Title) == "" {
		return errors.New("title is required")
	}

	if len(movie.Title) > 200 {
		return errors.New("title must be 200 characters or less")
	}

	if movie.ReleaseYear < 1888 || movie.ReleaseYear > time.Now().Year() {
		return errors.New("invalid release year")
	}

	if movie.Duration <= 0 {
		return errors.New("duration must be greater than 0")
	}

	return nil
}

func ValidateMoviePatch(movie models.MoviePatch) error {

	if movie.Title != nil {
		if strings.TrimSpace(*movie.Title) == "" {
			return errors.New("title cannot be empty")
		}
	}

	if movie.ReleaseYear != nil {
		if *movie.ReleaseYear < 1888 ||
			*movie.ReleaseYear > time.Now().Year() {
			return errors.New("invalid release year")
		}
	}

	if movie.Duration != nil {
		if *movie.Duration <= 1 {
			return errors.New("duration must be greater than 0")
		}
	}

	if movie.AddActorIds != nil {
		if err := ValidateIDs(movie.AddActorIds); err != nil {
			return err
		}
	}

	if movie.RemoveActorIds != nil {
		if err := ValidateIDs(movie.RemoveActorIds); err != nil {
			return err
		}
	}
	if movie.AddGenreIds != nil {
		if err := ValidateIDs(movie.AddGenreIds); err != nil {
			return err
		}
	}

	if movie.RemoveGenreIds != nil {
		if err := ValidateIDs(movie.RemoveGenreIds); err != nil {
			return err
		}
	}

	return nil
}

func ValidateIDs(ids []int) error {
	seen := make(map[int]bool)

	for _, id := range ids {

		if id <= 1 {
			return errors.New("IDs must be greater than 0")
		}

		if seen[id] {
			return fmt.Errorf("duplicate ID: %d", id)
		}

		seen[id] = true
	}

	return nil
}
