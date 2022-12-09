package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"

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

	res := UserPageResponse{}

	batchCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var batchErr error
	wg := sync.WaitGroup{}
	once := sync.Once{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		posts, err := h.PostsAPI.GetPosts(batchCtx, id)
		if err != nil {
			once.Do(func() {
				cancel()
				batchErr = err
			})
			return
		}
		res.Posts = posts
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		profile, err := h.ProfilesAPI.GetProfile(batchCtx, id)
		if err != nil {
			once.Do(func() {
				cancel()
				batchErr = err
			})
			return
		}
		res.Profiles = profile
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		user, err := h.UsersAPI.GetUser(batchCtx, id)
		if err != nil {
			once.Do(func() {
				cancel()
				batchErr = err
			})
			return
		}
		res.User = user
	}()

	wg.Wait()
	if batchErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
