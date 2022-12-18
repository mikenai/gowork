package handlers

import (
	"context"

	"github.com/mikenai/gowork/internal/models"
	"github.com/mikenai/gowork/internal/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CreateUserParams struct {
	Name string
}

//go:generate moq -rm -out users_mock.go . UsersService
type UsersService interface {
	GetOne(ctx context.Context, id string) (models.User, error)
}

type Users struct {
	protobuf.UnimplementedUsersServer
	user UsersService
}

func NewUsers(us UsersService) Users {
	return Users{user: us}
}

func (u *Users) GetUser(ctx context.Context, req *protobuf.GetRequest) (*protobuf.User, error) {
	user, err := u.user.GetOne(ctx, string(req.Id))
	if err != nil {
		return nil, status.Error(codes.NotFound, "user was not found")
	}

	return &protobuf.User{Id: user.ID, Name: user.Name}, nil
}
