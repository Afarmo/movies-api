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

type ActorRepo struct {
	db *sql.DB
}

func NewActorRepo(db *sql.DB) *ActorRepo {
	return &ActorRepo{db: db}
}

func (r *ActorRepo) InsertActor(ctx context.Context, actor *models.Actor) error {
	query := `
        INSERT INTO Actor (name, birthDate)
        VALUES (?, ?)
    `
	result, err := r.db.ExecContext(ctx, query,
		strings.TrimSpace(actor.Name),
		actor.BirthDate,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()

	if err != nil {
		return err
	}
	actor.ID = int(id)

	return nil
}

func (r *ActorRepo) UpdateActor(ctx context.Context, id int64, actor *models.ActorPatch) (*models.Actor, error) {
	query := "Update Actor SET "
	args := []any{}

	if actor.Name != nil {
		query += "name = ?, "
		args = append(args, actor.Name)
	}
	if actor.BirthDate != nil {
		query += "birthdate = ?, "
		args = append(args, actor.BirthDate)
	}

	query = strings.TrimSuffix(query, ", ")
	query += " WHERE id = ?"
	args = append(args, id)
	result, err := r.db.ExecContext(ctx, query, args...)

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

	a, err := r.ListOneActor(ctx, id)
	return a, nil
}

func (r *ActorRepo) DeleteActor(ctx context.Context, actorID int64, force bool) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int

	err = tx.QueryRowContext(ctx, `
		SELECT COUNT(*)
		FROM Movie_Actors
		WHERE actor_id = ?
	`, actorID).Scan(&count)

	if err != nil {
		return err
	}

	if count > 0 && !force {
		return fmt.Errorf("%w: actor is associated with %d movie(s)", apperrors.ErrAssociated, count)
	}

	if force {
		if _, err = tx.ExecContext(ctx, `
			DELETE FROM Movie_Actors
			WHERE actor_id = ?
		`, actorID); err != nil {
			return err
		}
	}

	result, err := tx.ExecContext(ctx, `
		DELETE FROM Actor
		WHERE id = ?
	`, actorID)

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

func (r *ActorRepo) ListAllActors(ctx context.Context) ([]*models.Actor, error) {
	query := "SELECT id, name, birthdate FROM Actor"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var actors []*models.Actor
	for rows.Next() {
		actor := &models.Actor{}
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		)
		if err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return actors, nil
}

func (r *ActorRepo) ListOneActor(ctx context.Context, id int64) (*models.Actor, error) {
	query := "SELECT * FROM Actor WHERE id = ?"

	actor := &models.Actor{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(&actor.ID, &actor.Name, &actor.BirthDate)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return actor, nil
}

func (r *ActorRepo) CountActors(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Actor").Scan(&count)
	return count, err
}

func (r *ActorRepo) SearchActors(ctx context.Context, name string) ([]*models.Actor, error) {
	query := "SELECT * FROM Actor WHERE lower(name) LIKE ?"

	rows, err := r.db.QueryContext(ctx, query, "%"+strings.TrimSpace(name)+"%")
	if err != nil {
		return nil, err
	}

	var actors []*models.Actor

	for rows.Next() {
		actor := &models.Actor{}

		if err := rows.Scan(&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		); err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}

func (r *ActorRepo) ActorsByMovie(ctx context.Context, movieId int64) ([]*models.Actor, error) {
	var exists int

	err := r.db.QueryRowContext(
		ctx,
		`SELECT 1 FROM Movie WHERE id = ?`,
		movieId,
	).Scan(&exists)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	query := `
        SELECT a.id, a.name, a.birthDate
        FROM Actor a
        JOIN Movie_Actors ma ON a.id = ma.actor_id
        WHERE ma.movie_id = ?
    `

	rows, err := r.db.QueryContext(ctx, query, movieId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var actors []*models.Actor

	for rows.Next() {
		actor := &models.Actor{}
		if err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&actor.BirthDate,
		); err != nil {
			return nil, err
		}

		actors = append(actors, actor)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return actors, nil
}
