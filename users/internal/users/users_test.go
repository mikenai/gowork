package users

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mikenai/gowork/internal/models"
)

func TestService_Create(t *testing.T) {
	type fields struct {
		repo Repositry
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
	}{
		{
			name:    "invelid name: empty",
			fields:  fields{},
			args:    args{},
			want:    models.User{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}
			got, err := s.Create(context.Background(), tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("Service.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	type fields struct {
		repo Repositry
	}
	type args struct {
		ctx    context.Context
		params UserUpdateParams
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    models.User
		wantErr bool
		errorIs error
	}{
		{
			name:   "invalid argument: empty name",
			fields: fields{},
			args: args{
				ctx: nil,
				params: UserUpdateParams{
					Name: "",
				},
			},
			want:    models.User{},
			wantErr: true,
			errorIs: models.InvalidArgumentErr,
		},
		{
			name:   "invalid argument: len restiction",
			fields: fields{},
			args: args{
				ctx: nil,
				params: UserUpdateParams{
					Name: "12345678901",
				},
			},
			want:    models.User{},
			wantErr: true,
			errorIs: models.InvalidArgumentErr,
		},
		{
			name: "repository returns an error",
			fields: fields{
				repo: &RepositryMock{
					UpdateFunc: func(ctx context.Context, params UserUpdateParams) (models.User, error) {
						return models.User{}, errors.New("some error")
					},
				},
			},
			args: args{
				ctx: nil,
				params: UserUpdateParams{
					Name: "Gabriella",
				},
			},
			want:    models.User{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				repo: tt.fields.repo,
			}

			got, err := s.Update(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr || (tt.errorIs != nil && !errors.Is(err, tt.errorIs)) {
				t.Errorf("Service.Update() error = %v, wantErr %v want error as %v", err, tt.wantErr, tt.errorIs)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Service.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}
