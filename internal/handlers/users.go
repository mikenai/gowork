package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikenai/gowork/internal/models"
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
	return Users{}
}

func (u Users) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", u.Create)

	return r
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var userParams CreateUserParams
	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := u.user.Create(ctx, userParams.Name)
	if err != nil {
		if errors.Is(err, models.UserCreateParamInvalidNameErr) {
			http.Error(w, "", http.StatusBadRequest)
		}
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
