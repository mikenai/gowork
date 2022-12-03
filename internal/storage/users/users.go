package users

import (
	"context"
	"database/sql"
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
