package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikenai/gowork/internal/models"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) Storage {
	return Storage{db: db}
}

func (s Storage) Create(ctx context.Context, name string) (models.User, error) {
	id := uuid.NewString()
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to execute insert: %w", err)
	}
	return models.User{ID: id, Name: name}, nil
}

func (s Storage) GetByID(ctx context.Context, id string) (models.User, error) {
	var usr models.User

	row := s.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id=?", id)

	if err := row.Scan(&usr.ID, &usr.Name); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("failed to get user: %w", models.NotFoundErr)
		}
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return usr, nil
}
