package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"movies-api/internal/apperrors"
	"movies-api/internal/models"
	"strings"
)

type GenreRepo struct {
	db *sql.DB
}

func NewGenreRepo(db *sql.DB) *GenreRepo {
	return &GenreRepo{db: db}
}

func (r *GenreRepo) InsertGenre(ctx context.Context, genre *models.Genre) error {
	query := "INSERT INTO Genre (name) VALUES(?)"

	result, err := r.db.ExecContext(ctx, query, genre.Name)
	if err != nil {
		if isUniqueConstraintError(err) {
			return apperrors.ErrConflict
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil
	}
	genre.ID = int(id)

	return nil
}

func (r *GenreRepo) UpdateGenre(ctx context.Context, genre *models.Genre) error {
	query := "UPDATE Genre SET name = ? WHERE id = ?"

	result, err := r.db.ExecContext(ctx, query, genre.Name, genre.ID)
	if err != nil {
		if isUniqueConstraintError(err) {
			return apperrors.ErrConflict
		}
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (r *GenreRepo) DeleteGenre(ctx context.Context, genreID int64, force bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int

	if err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM Movie_Genres
		WHERE genre_id = ?
	`, genreID).Scan(&count); err != nil {
		return err
	}

	if count > 0 && !force {
		return fmt.Errorf("%w: genre is associated with %d movie(s)", apperrors.ErrAssociated, count)
	}

	if force {
		if _, err = tx.ExecContext(ctx, `
			DELETE FROM Movie_Genres
			WHERE genre_id = ?
		`, genreID); err != nil {
			return err
		}
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM Genre
		WHERE id = ?
	`, genreID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return tx.Commit()

}

func (r *GenreRepo) ListGenres(ctx context.Context, filter *models.GenreFilter) (*models.ListGenresResult, error) {
	var total int

	if err := r.db.QueryRowContext(
		ctx,
		"SELECT COUNT(*) FROM Genre",
	).Scan(&total); err != nil {
		return nil, err
	}

	query := "SELECT id, name FROM Genre"

	var args []any

	if filter.Page != nil && filter.Size != nil {
		offset := *filter.Page * *filter.Size
		query += " LIMIT ? OFFSET ?"
		args = append(args, *filter.Size, offset)
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var genres []*models.Genre
	for rows.Next() {
		genre := &models.Genre{}
		if err := rows.Scan(&genre.ID, &genre.Name); err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.ListGenresResult{
		Genres: genres,
		Total:  total,
	}, nil
}

func (r *GenreRepo) ListOneGenre(ctx context.Context, id int64) (*models.Genre, error) {
	query := "SELECT * FROM Genre WHERE id = ?"

	genre := &models.Genre{}
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&genre.ID, &genre.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return genre, nil
}

func isUniqueConstraintError(err error) bool {
	// SQLite error message contains "UNIQUE constraint failed"
	return err != nil && (strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "unique constraint"))
}
