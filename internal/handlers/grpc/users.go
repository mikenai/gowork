package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/mikenai/gowork/internal/models"
	pb "github.com/mikenai/gowork/internal/proto"
	"github.com/mikenai/gowork/internal/users"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UsersGRPC struct {
	pb.UnimplementedUsersServiceServer
	users users.Service
	log   zerolog.Logger
}

func (grpc *UsersGRPC) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error) {
	usr, err := grpc.users.GetOne(ctx, in.Id)
	if err != nil {
		grpc.log.Error().Err(err).Msgf("Fetch failed: %v", err)
		msg := fmt.Sprintf("Failed to get user by id: %s", in.Id)
		if errors.Is(err, models.NotFoundErr) {
			return &pb.UserResponse{}, status.Error(codes.NotFound, msg)
		}

		return &pb.UserResponse{}, status.Error(codes.Internal, msg)
	}

	grpc.log.Info().Msgf("Fetch %s: %s", usr.ID, usr.Name)
	return &pb.UserResponse{Id: usr.ID, Name: usr.Name}, nil
}

func NewUsersGRPC(users users.Service, log zerolog.Logger) *UsersGRPC {
	return &UsersGRPC{
		users: users,
		log:   log,
	}
}
