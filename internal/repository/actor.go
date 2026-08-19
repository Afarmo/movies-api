package repository

import (
	"context"
	"database/sql"
	"errors"
	"movies-api/internal/models"
	"strings"
	"time"
)

var (
	ErrNotFound     = errors.New("record not found")
	ErrDuplicateID  = errors.New("duplicate key violation")
	ErrInvalidInput = errors.New("invalid input")
)

// ActorRepo handles user database operations
type ActorRepo struct {
	db *sql.DB
}

// NewActorRepo creates new actor repo
func NewActorRepo(db *sql.DB) *ActorRepo {
	return &ActorRepo{db: db}
}

// Creates and Inserts a new actor into the db
// Returns the insterted Actor
func (r *ActorRepo) InsertActor(ctx context.Context, actor *models.Actor) error {

	query := `
        INSERT INTO Actor (name, birthDate)
        VALUES (?, ?)
    `
	result, err := r.db.ExecContext(ctx, query,
		strings.TrimSpace(actor.Name),
		actor.BirthDate.Format("2006-01-02"),
	)
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
	actor.ID = int(id)

	return nil
}

// Update an Actor
func (r *ActorRepo) UpdateActor(ctx context.Context, actor *models.Actor) error {
	query := `
        UPDATE Actor
        SET name = ?, birthDay = ?
        WHERE id = ?
    `
	result, err := r.db.ExecContext(ctx, query,
		actor.Name,
		actor.BirthDate,
		actor.ID,
	)

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

// Delete an Actor from the db
func (r *ActorRepo) DeleteActor(ctx context.Context, id int64) error {
	query := `DELETE FROM Actor WHERE id = ?`

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

func (r *ActorRepo) ListAllActors(ctx context.Context) ([]*models.Actor, error) {
	query := "SELECT * FROM Actor"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	var actors []*models.Actor
	var bd string
	for rows.Next() {
		actor := &models.Actor{}
		err := rows.Scan(
			&actor.ID,
			&actor.Name,
			&bd,
		)
		if err != nil {
			return nil, err
		}
		actor.BirthDate, _ = time.Parse("2006-01-02", bd)

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
	var bd string
	err := r.db.QueryRowContext(ctx, query, id).Scan(&actor.ID, &actor.Name, &bd)
	actor.BirthDate, _ = time.Parse("2006-01-02", bd)
	if err != nil {
		return nil, err
	}

	return actor, nil
}

func (r *ActorRepo) CountActors(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM Actor").Scan(&count)
	return count, err
}

func (r *ActorRepo) ActorsByName(ctx context.Context, name string) ([]*models.Actor, error) {
	query := "SELECT * FROM Actor WHERE name = ?"
	var bd string

	rows, err := r.db.QueryContext(ctx, query, strings.TrimSpace(name))
	if err != nil {
		return nil, err
	}

	var actors []*models.Actor

	for rows.Next() {

		actor := &models.Actor{}
		err := rows.Scan(&actor.ID,
			&actor.Name,
			&bd)
		actor.BirthDate, _ = time.Parse("2006-01-02", bd)
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

func isUniqueConstraintError(err error) bool {
	// SQLite error message contains "UNIQUE constraint failed"
	return err != nil && (strings.Contains(err.Error(), "UNIQUE constraint failed") ||
		strings.Contains(err.Error(), "unique constraint"))
}
