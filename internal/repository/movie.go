package repository

import (
	"context"
	"database/sql"
	"movies-api/internal/models"
)

type MovieRepo struct {
	db *sql.DB
}

func NewMovieRepo(db *sql.DB) *MovieRepo {
	return &MovieRepo{db: db}
}

func (r *MovieRepo) InsertMovie(ctx context.Context, movie *models.Movie) error {

	query := `
        INSERT INTO Movie (title, releaseYear, duration)
        VALUES (?, ?, ?)
    `
	result, err := r.db.ExecContext(ctx, query, movie.Title, movie.ReleaseYear, movie.Duration)
	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrDuplicateID
		}
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	movie.ID = int(id)
	return nil
}

func (r *MovieRepo) UpdateMovie(ctx context.Context, movie *models.Movie) error {

	query := ` Update Movie
		SET title = ?, releaseYear = ?, duration = ?
		where id = ?
	`
	result, err := r.db.ExecContext(ctx, query, movie.Title, movie.ReleaseYear, movie.Duration, movie.ID)

	if err != nil {
		if isUniqueConstraintError(err) {
			return ErrDuplicateID
		}
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

func (r *MovieRepo) DeleteMovie(ctx context.Context, id int64) error {

	query := "DELETE FROM Movie WHERE id = ?"
	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *MovieRepo) ListOneMovie(ctx context.Context, id int64) (*models.Movie, error) {

	movie := &models.Movie{}
	query := "SELECT * FROM Movie WHERE id = ?"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&movie.ID,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration)
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (r *MovieRepo) ListAllMovies(ctx context.Context) ([]*models.Movie, error) {
	query := "SELECT * FROM Movie"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil

}

func (r *MovieRepo) MovieByRealeaseYear(ctx context.Context, year string) ([]*models.Movie, error) {
	query := "SELECT * FROM Movie WHERE releaseYear = ? "
	rows, err := r.db.QueryContext(ctx, query, year)
	if err != nil {
		return nil, err
	}
	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}
func (r *MovieRepo) MoviesByGenre(ctx context.Context, genreId int64) ([]*models.Movie, error) {
	query := `
        SELECT m.id, m.title, m.releaseYear, m.duration
        FROM Movie m
        JOIN Movie_Genres mg ON m.id = mg.movie_id
        WHERE mg.genre_id = ?
    `
	rows, err := r.db.QueryContext(ctx, query, genreId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (r *MovieRepo) MoviesByActor(ctx context.Context, actorId int64) ([]*models.Movie, error) {
	query := `
        SELECT m.id, m.title, m.releaseYear, m.duration
        FROM Movie m
        JOIN Movie_Actors ma ON m.id = ma.movie_id
        WHERE ma.actor_id = ?
    `
	rows, err := r.db.QueryContext(ctx, query, actorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var movies []*models.Movie
	for rows.Next() {
		movie := &models.Movie{}
		err := rows.Scan(&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration)
		if err != nil {
			return nil, err
		}
		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}
func (r *MovieRepo) MovieByName(ctx context.Context, name string) (*models.Movie, error) {
	query := "SELECT * FROM Movie WHERE name = ?"
	movie := &models.Movie{}
	err := r.db.QueryRowContext(ctx, query, name).Scan(&movie.ID,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration)
	if err != nil {
		return nil, err
	}
	return movie, nil
}
