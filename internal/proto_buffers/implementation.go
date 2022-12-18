package proto_buffers

import (
	"context"
	"errors"

	"github.com/mikenai/gowork/internal/models"
	"github.com/mikenai/gowork/internal/users"
	"github.com/rs/zerolog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type implementation struct {
	UnimplementedUsersServiceServer
	users users.Service
	log   zerolog.Logger
}

func (s *implementation) GetOne(ctx context.Context, in *GetUserRequest) (*UserResponse, error) {
	usr, err := s.users.GetOne(ctx, in.Id)
	if err != nil {
		s.log.Error().Err(err).Msgf("Fetch failed: %v", err)

		if errors.Is(err, models.NotFoundErr) {
			return &UserResponse{}, status.Errorf(codes.NotFound, "Failed to get user")
		}

		return &UserResponse{}, status.Errorf(codes.Internal, "Failed to get user")
	}

	s.log.Info().Msgf("Fetch %s: %s", usr.ID, usr.Name)
	return &UserResponse{Id: usr.ID, Name: usr.Name}, nil
}

func Implementation(users users.Service, log zerolog.Logger) *implementation {
	return &implementation{
		users: users,
		log:   log,
	}
}
