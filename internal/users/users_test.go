package users

import (
	"context"
	"testing"

	"github.com/mikenai/gowork/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestService_Create(t *testing.T) {

	type fields struct {
		repo Repository
	}
	type args struct {
		name string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      models.User
		wantErrIs error
	}{
		{
			name: "success",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{ID: "1", Name: name}, nil
					},
				},
			},
			args: args{
				name: "Rob",
			},
			want:      models.User{ID: "1", Name: "Rob"},
			wantErrIs: nil,
		},
		{
			name: "empty name",
			fields: fields{
				repo: &RepositoryMock{},
			},
			args: args{
				name: "",
			},
			want:      models.User{},
			wantErrIs: models.UserCreateParamInvalidNameErr,
		},
		{
			name: "saving error",
			fields: fields{
				repo: &RepositoryMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{}, models.InvalidErr
					},
				},
			},
			args: args{
				name: "test",
			},
			want:      models.User{},
			wantErrIs: models.InvalidErr,
		},
	}
	for _, tt := range tests {
		ctx := context.Background()
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			got, err := s.Create(ctx, tt.args.name)
			assert.Equal(t, tt.want, got)
			assert.ErrorIs(t, err, tt.wantErrIs)
		})
	}
}
