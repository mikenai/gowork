package usersapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/mikenai/gowork/internal/shared/protobuf"
	"google.golang.org/grpc"
)

type User struct {
	ID   string
	Name string
}

type Client struct {
	BaseURL string
	Http    http.Client
}

func (cl *Client) GetUser(ctx context.Context, id string) (User, error) {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(cl.BaseURL, grpc.WithInsecure())
	if err != nil {
		errors.New("connection issue")
	}
	defer conn.Close()

	client := protobuf.NewUsersClient(conn)

	req := protobuf.GetRequest{Id: id}
	user, err := client.GetUser(ctx, &req)
	if err != nil {
		return User{}, fmt.Errorf("json error: %w", err)
	}

	return User{ID: user.Id, Name: user.Name}, nil

}

