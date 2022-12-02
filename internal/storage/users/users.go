package users

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/mikenai/gowork/internal/models"
)

const file string = "users.db"

type Storage struct {
	db *sql.DB
}

func NewStorage() (*Storage, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, err
	}

	return &Storage{db}, nil
}

func (s Storage) Create(ctx context.Context, name string) (models.User, error) {
	id := uuid.NewString()
	_, err := s.db.ExecContext(ctx, "INSERT INTO users (id, name) VALUES (?, ?)", id, name)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to execute insert: %w", err)
	}
	return models.User{ID: id, Name: name}, nil
}

func (s Storage) Fetch(ctx context.Context, id string) (models.User, error) {
	row := s.db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id=?", id)

	user := models.User{}
	if err := row.Scan(&user.ID, &user.Name); err != nil {
		return models.User{}, err
	}

	return user, nil
}
