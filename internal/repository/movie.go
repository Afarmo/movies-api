package repository

import (
	"context"
	"database/sql"
	"movies-api/internal/models"
	"strings"
)

type MovieRepo struct {
	db *sql.DB
}

func NewMovieRepo(db *sql.DB) *MovieRepo {
	return &MovieRepo{db: db}
}

func (r *MovieRepo) InsertMovie(ctx context.Context, movie *models.Movie) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer tx.Rollback()
	query := `
        INSERT INTO Movie (title, releaseYear, duration)
        VALUES (?, ?, ?)
    `
	result, err := r.db.ExecContext(ctx, query, movie.Title, movie.ReleaseYear, movie.Duration)
	if err != nil {
		return err
	}
	movieID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	movie.ID = int(movieID)
	for _, actorID := range movie.ActorIds {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO Movie_Actors (movie_id, actor_id)
             VALUES (?, ?)`,
			movieID,
			actorID,
		)

		if err != nil {
			return err
		}
	}
	for _, genreID := range movie.GenreIds {
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO Movie_Genres (movie_id, genre_id)
             VALUES (?, ?)`,
			movieID,
			genreID,
		)

		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *MovieRepo) UpdateMovie(ctx context.Context, id int64, movie *models.MoviePatch) error {

	query := " Update Movie SET "
	args := []any{}

	if movie.Title != nil {
		query += "title = ?, "
		args = append(args, movie.Title)
	}

	if movie.ReleaseYear != nil {
		query += "releaseyear = ?, "
		args = append(args, movie.ReleaseYear)
	}

	if movie.Duration != nil {
		query += "duration = ?, "
		args = append(args, movie.Duration)
	}

	query = strings.TrimSuffix(query, ", ")

	query += " WHERE id = ?"
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, query, args...)
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
	query := "SELECT id, title, releaseyear, duration FROM Movie WHERE id = ?"
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
	query := "SELECT id, title, releaseyear, duration FROM Movie"
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

func (r *MovieRepo) SearchMovies(ctx context.Context, name string) ([]*models.Movie, error) {
	query := "SELECT id, title, releaseyear, duration FROM Movie WHERE title = ?"
	rows, err := r.db.QueryContext(ctx, query, name)
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

func (r *MovieRepo) MovieByRealeaseYear(ctx context.Context, year int) ([]*models.Movie, error) {
	query := "SELECT id, title, releaseyear, duration FROM Movie WHERE releaseYear = ? "
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
