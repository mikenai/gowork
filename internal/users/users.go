package users

import (
	"context"
	"fmt"

	"github.com/mikenai/gowork/internal/models"
)

//go:generate moq -rm -out users_repository_mock.go . Repository
type Repository interface {
	Create(ctx context.Context, name string) (models.User, error)
	Fetch(ctx context.Context, id string) (models.User, error)
}

type Service struct {
	repo Repository
}

func (s Service) Create(ctx context.Context, name string) (models.User, error) {
	if name == "" {
		return models.User{}, fmt.Errorf("invalid name argument: %w", models.UserCreateParamInvalidNameErr)
	}

	usr, err := s.repo.Create(ctx, name)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return usr, nil
}

func (s Service) Fetch(ctx context.Context, id string) (models.User, error) {
	if id == "" {
		return models.User{}, fmt.Errorf("not found: %w", models.NotFound)
	}

	usr, err := s.repo.Fetch(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("not found: %w", err)
	}

	return usr, nil
}
