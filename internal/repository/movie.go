package repository

import (
	"context"
	"database/sql"
	"errors"
	"movies-api/internal/apperrors"
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
	result, err := tx.ExecContext(ctx, query, movie.Title, movie.ReleaseYear, movie.Duration)
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

func (r *MovieRepo) UpdateMovie(ctx context.Context, id int64, movie *models.MoviePatch) (*models.Movie, error) {
	tx, err := r.db.BeginTx(ctx, nil)

	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	query := " Update Movie SET "
	var args []any

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
	if args != nil {
		args = append(args, id)

		result, err := tx.ExecContext(ctx, query, args...)
		if err != nil {
			return nil, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return nil, err
		}

		if rowsAffected == 0 {
			return nil, apperrors.ErrNotFound
		}
	}

	if movie.AddActorIds != nil {
		for _, actorID := range movie.AddActorIds {
			_, err := tx.ExecContext(ctx, `
		INSERT INTO Movie_Actors (movie_id, actor_id)
		VALUES (?, ?)
	`, id, actorID)

			if err != nil {
				return nil, apperrors.ErrDuplicateID
			}
		}

	}

	if movie.RemoveActorIds != nil {
		for _, actorID := range movie.RemoveActorIds {
			_, err := tx.ExecContext(ctx, `
		DELETE FROM Movie_Actors
		WHERE movie_id = ? AND actor_id = ?
	`, id, actorID)

			if err != nil {

				return nil, err
			}
		}
	}

	if movie.AddGenreIds != nil {
		for _, genreID := range movie.AddGenreIds {
			_, err := tx.ExecContext(ctx, `
		INSERT INTO Movie_Genres (movie_id, genre_id)
		VALUES (?, ?)
	`, id, genreID)

			if err != nil {
				return nil, err
			}
		}
	}

	if movie.RemoveGenreIds != nil {
		for _, genreID := range movie.RemoveGenreIds {
			_, err := tx.ExecContext(ctx, `
		DELETE FROM Movie_Genres
		WHERE movie_id = ? AND genre_id = ?
	`, id, genreID)

			if err != nil {
				return nil, err
			}
		}
	}
	tx.Commit()
	m, err := r.ListOneMovie(ctx, id)
	if err != nil {
		return nil, err
	}

	return m, nil
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
		return apperrors.ErrNotFound
	}
	return nil
}

func (r *MovieRepo) ListOneMovie(ctx context.Context, id int64) (*models.Movie, error) {

	movie := &models.Movie{}
	query := "SELECT id, title, releaseyear, duration FROM Movie WHERE id = ?"
	if err := r.db.QueryRowContext(ctx, query, id).Scan(&movie.ID,
		&movie.Title,
		&movie.ReleaseYear,
		&movie.Duration,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	actorQuery := `
			SELECT actor_id
			FROM Movie_Actors
			WHERE movie_id = ?
		`

	actorRows, err := r.db.QueryContext(ctx, actorQuery, movie.ID)
	if err != nil {
		return nil, err
	}

	for actorRows.Next() {
		var actorID int

		err := actorRows.Scan(&actorID)
		if err != nil {
			actorRows.Close()
			return nil, err
		}

		movie.ActorIds = append(movie.ActorIds, actorID)
	}

	if err := actorRows.Err(); err != nil {
		actorRows.Close()
		return nil, err
	}

	actorRows.Close()

	genreQuery := `
			SELECT genre_id
			FROM Movie_Genres
			WHERE movie_id = ?
		`

	genreRows, err := r.db.QueryContext(ctx, genreQuery, movie.ID)
	if err != nil {
		return nil, err
	}

	for genreRows.Next() {
		var genreID int

		err := genreRows.Scan(&genreID)
		if err != nil {
			genreRows.Close()
			return nil, err
		}

		movie.GenreIds = append(movie.GenreIds, genreID)
	}

	if err := genreRows.Err(); err != nil {
		genreRows.Close()
		return nil, err
	}

	genreRows.Close()
	return movie, nil
}

func (r *MovieRepo) ListMovies(ctx context.Context, filter *models.MovieFilter) ([]*models.Movie, error) {
	query := `
		SELECT DISTINCT m.id, m.title, m.releaseYear, m.duration
		FROM Movie m
	`

	var conditions []string
	var args []any

	if filter.GenreID != nil {
		query += `
		JOIN Movie_Genres mg ON m.id = mg.movie_id
	`
		conditions = append(conditions, "mg.genre_id = ?")
		args = append(args, *filter.GenreID)
	}

	if filter.ActorID != nil {
		query += `
		JOIN Movie_Actors ma ON m.id = ma.movie_id
	`
		conditions = append(conditions, "ma.actor_id = ?")
		args = append(args, *filter.ActorID)
	}

	if filter.Year != nil {
		conditions = append(conditions, "m.releaseYear = ?")
		args = append(args, *filter.Year)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

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

	var movies []*models.Movie

	for rows.Next() {
		movie := &models.Movie{}

		if err := rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseYear,
			&movie.Duration,
		); err != nil {
			return nil, err
		}

		actorQuery := `
		SELECT actor_id
		FROM Movie_Actors
		WHERE movie_id = ?
		`
		actorRows, err := r.db.QueryContext(ctx, actorQuery, movie.ID)
		if err != nil {
			return nil, err
		}
		for actorRows.Next() {
			var actorID int
			if err := actorRows.Scan(&actorID); err != nil {
				actorRows.Close()
				return nil, err
			}

			movie.ActorIds = append(movie.ActorIds, actorID)
		}

		if err := actorRows.Err(); err != nil {
			actorRows.Close()
			return nil, err
		}
		actorRows.Close()

		genreQuery := `
		SELECT genre_id
		FROM Movie_Genres
		WHERE movie_id = ?
		`
		genreRows, err := r.db.QueryContext(ctx, genreQuery, movie.ID)
		if err != nil {
			return nil, err
		}
		for genreRows.Next() {
			var genreID int
			if err := genreRows.Scan(&genreID); err != nil {
				genreRows.Close()
				return nil, err
			}

			movie.GenreIds = append(movie.GenreIds, genreID)
		}

		if err := genreRows.Err(); err != nil {
			genreRows.Close()
			return nil, err
		}
		genreRows.Close()

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepo) SearchMovies(ctx context.Context, name string) ([]*models.Movie, error) {
	query := "SELECT id, title, releaseyear, duration FROM Movie WHERE lower(title) LIKE ?"
	rows, err := r.db.QueryContext(ctx, query, "%"+strings.TrimSpace(name)+"%")
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
		actorQuery := `
			SELECT actor_id
			FROM Movie_Actors
			WHERE movie_id = ?
		`

		actorRows, err := r.db.QueryContext(ctx, actorQuery, movie.ID)
		if err != nil {
			return nil, err
		}

		for actorRows.Next() {
			var actorID int

			err := actorRows.Scan(&actorID)
			if err != nil {
				actorRows.Close()
				return nil, err
			}

			movie.ActorIds = append(movie.ActorIds, actorID)
		}

		if err := actorRows.Err(); err != nil {
			actorRows.Close()
			return nil, err
		}

		actorRows.Close()

		genreQuery := `
			SELECT genre_id
			FROM Movie_Genres
			WHERE movie_id = ?
		`

		genreRows, err := r.db.QueryContext(ctx, genreQuery, movie.ID)
		if err != nil {
			return nil, err
		}

		for genreRows.Next() {
			var genreID int

			err := genreRows.Scan(&genreID)
			if err != nil {
				genreRows.Close()
				return nil, err
			}

			movie.GenreIds = append(movie.GenreIds, genreID)
		}

		if err := genreRows.Err(); err != nil {
			genreRows.Close()
			return nil, err
		}

		genreRows.Close()

		movies = append(movies, movie)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}
