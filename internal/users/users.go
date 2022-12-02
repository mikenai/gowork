package users

import (
	"context"
	"fmt"

	"github.com/mikenai/gowork/internal/models"
)

//go:generate moq -rm -out users_repository_mock.go . Repository
type Repository interface {
	Create(ctx context.Context, name string) (models.User, error)
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
