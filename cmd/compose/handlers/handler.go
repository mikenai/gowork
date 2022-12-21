package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/mikenai/gowork/cmd/compose/pkg/stub"
	"github.com/mikenai/gowork/cmd/compose/pkg/usersapi"
	pb "github.com/mikenai/gowork/internal/proto"
	"github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type Posts interface {
	GetPosts(ctx context.Context, id string) ([]stub.Post, error)
}

type Profiles interface {
	GetProfile(ctx context.Context, id string) (stub.Profile, error)
}

type Users interface {
	GetUser(ctx context.Context, in *pb.GetUserRequest, opts ...grpc.CallOption) (*pb.UserResponse, error)
}

type Handler struct {
	PostsAPI    Posts
	ProfilesAPI Profiles
	UsersGRPC   Users

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

	eg, batchCtx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		posts, err := h.PostsAPI.GetPosts(batchCtx, id)
		if err != nil {
			return err
		}
		res.Posts = posts
		return nil
	})

	eg.Go(func() error {
		profile, err := h.ProfilesAPI.GetProfile(batchCtx, id)
		if err != nil {
			return err
		}
		res.Profiles = profile
		return nil
	})

	eg.Go(func() error {
		user, err := h.UsersGRPC.GetUser(batchCtx, &pb.GetUserRequest{Id: id})
		if err != nil {
			return err
		}
		res.User = usersapi.User{ID: user.Id, Name: user.Name}
		return nil
	})

	batchErr := eg.Wait()
	if batchErr != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(res)
}
