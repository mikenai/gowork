package handlers

import (
	"encoding/json"
	"net/http"
)

type CreateUserParams struct {
	Name string
}

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

//go:generate moq -rm -out users_mock.go . UsersService
type UsersService interface {
	Create(name string) (User, error)
}

type Users struct {
	user UsersService
}

func (u Users) Create(w http.ResponseWriter, r *http.Request) {
	var userParams CreateUserParams

	if err := json.NewDecoder(r.Body).Decode(&userParams); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	user, err := u.user.Create(userParams.Name)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}
}
