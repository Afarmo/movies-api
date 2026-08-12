package repository

import (
	"context"
	"database/sql"
	"movies-api/internal/models"
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
			return ErrDuplicateId
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
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrNotFound
	}
	return nil

}

func (r *GenreRepo) DeleteGenre(ctx context.Context, id int64) error {
	query := "DELETE FROM Genre WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	ra, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if ra == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *GenreRepo) ListAllGenres(ctx context.Context) ([]*models.Genre, error) {
	query := "SELECT * FROM Genre"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var genres []*models.Genre
	for rows.Next() {
		genre := &models.Genre{}
		err := rows.Scan(&genre.ID, &genre.Name)
		if err != nil {
			return nil, err
		}
		genres = append(genres, genre)
	}
	return genres, nil
}

func (r *GenreRepo) ListOneGenre(ctx context.Context, id int64) (*models.Genre, error) {
	query := "SELECT * FROM Genre WHERE id = ?"
	genre := &models.Genre{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&genre.ID, &genre.Name)
	if err != nil {
		return nil, err
	}
	return genre, nil
}
