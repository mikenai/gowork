package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikenai/gowork/cmd/compose/pkg/stub"
	"github.com/mikenai/gowork/cmd/compose/pkg/usersapi"
	"github.com/rs/zerolog"
)

type Posts interface {
	GetPosts(ctx context.Context, id string) ([]stub.Post, error)
}

type Profiles interface {
	GetProfile(ctx context.Context, id string) (stub.Profile, error)
}

type Users interface {
	GetUser(ctx context.Context, id string) (usersapi.User, error)
}

type Handler struct {
	PostsAPI    Posts
	ProfilesAPI Profiles
	UsersAPI    Users

	Log zerolog.Logger
}

type UserPageResponse struct {
	Posts    []stub.Post
	User     usersapi.User
	Profiles stub.Profile
}

func (h Handler) UserPage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "user_id")
	ctx := r.Context()

	posts, err := h.PostsAPI.GetPosts(ctx, id)
	if err != nil {
		h.Log.Error().Err(err).Msg("get posts error")
		http.Error(w, "get posts", http.StatusInternalServerError)
		return
	}

	profile, err := h.ProfilesAPI.GetProfile(ctx, id)
	if err != nil {
		h.Log.Error().Err(err).Msg("get profile error")
		http.Error(w, "get profile", http.StatusInternalServerError)
		return
	}

	user, err := h.UsersAPI.GetUser(ctx, id)
	if err != nil {
		h.Log.Error().Err(err).Msg("get user error")
		http.Error(w, "get user", http.StatusInternalServerError)
		return
	}

	res := UserPageResponse{
		Posts:    posts,
		User:     user,
		Profiles: profile,
	}

	json.NewEncoder(w).Encode(res)
}
