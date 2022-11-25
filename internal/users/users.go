package users

import (
	"context"
	"fmt"

	"github.com/mikenai/gowork/internal/models"
)

type Repositry interface {
	Create(ctx context.Context, name string) (models.User, error)
}

type Service struct {
	repo Repositry
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
