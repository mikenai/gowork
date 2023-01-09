package usersapi

import (
	"context"
	"fmt"

	"github.com/mikenai/gowork/internal/shared/protobuf"
)

type User struct {
	ID   string
	Name string
}

type Client struct {
	Grpc protobuf.UsersClient
}

func (cl *Client) GetUser(ctx context.Context, id string) (User, error) {
	req := protobuf.GetRequest{Id: id}
	user, err := cl.Grpc.GetUser(ctx, &req)
	if err != nil {
		return User{}, fmt.Errorf("json error: %w", err)
	}

	return User{ID: user.Id, Name: user.Name}, nil
}
