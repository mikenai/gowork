package users

import (
	"context"
	"fmt"

	"github.com/mikenai/gowork/internal/models"
)

//go:generate moq -rm -out users_mock.go . Repositry
type Repositry interface {
	Create(ctx context.Context, name string) (models.User, error)
	GetByID(ctx context.Context, id string) (models.User, error)
	Update(ctx context.Context, params UserUpdateParams) (models.User, error)
}

type Service struct {
	repo Repositry
}

func New(r Repositry) Service {
	return Service{repo: r}
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

func (s Service) GetOne(ctx context.Context, id string) (models.User, error) {
	if id == "" {
		return models.User{}, fmt.Errorf("id is empty: %w", models.InvalidArgumentErr)
	}

	usr, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return usr, nil
}

type UserUpdateParams struct {
	Name string `json:"name"`
}

func (s Service) Update(ctx context.Context, params UserUpdateParams) (models.User, error) {
	if params.Name == "" {
		return models.User{}, models.UserUpdateParamInvalidNameEmptyErr
	}

	if len(params.Name) > 10 {
		return models.User{}, models.UserUpdateParamInvalidNameLenErr
	}

	user, err := s.repo.Update(ctx, params)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
