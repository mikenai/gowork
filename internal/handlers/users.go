package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikenai/gowork/internal/models"
	"github.com/mikenai/gowork/pkg/logger"
	"github.com/mikenai/gowork/pkg/response"
)

type CreateUserParams struct {
	Name string
}

//go:generate moq -rm -out users_mock.go . UsersService
type UsersService interface {
	Create(ctx context.Context, name string) (models.User, error)
}

type Users struct {
	user UsersService
}

func NewUsers(us UsersService) Users {
	return Users{user: us}
}

func (u Users) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", u.Create)

	return r
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logger.FromContext(ctx)

	var userParams CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		log.Error().Err(err).Msg("failed to parse params")
		response.InternalError(w)
		return
	}

	user, err := u.user.Create(ctx, userParams.Name)
	if err != nil {
		if errors.Is(err, models.UserCreateParamInvalidNameErr) {
			response.BadRequest(w)
		}
		log.Error().Err(err).Msg("failed to create user")
		response.InternalError(w)
		return
	}

	if err := response.JSON(w, user); err != nil {
		log.Error().Err(err).Msg("failed to encode response")
	}
}
