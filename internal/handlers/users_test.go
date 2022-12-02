package handlers

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mikenai/gowork/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUsers_Create(t *testing.T) {
	type fields struct {
		user UsersService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody []byte
	}{
		{
			name: "success",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						assert.Equal(t, "mike", name)
						return models.User{Name: name, ID: "1"}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusOK,
			wantBody: []byte(`{"id":"1","name":"mike"}` + "\n"),
		},
		{
			name: "cant save",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{}, fmt.Errorf("Can't create user")
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": "mike"}`)),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: []byte{0xa},
		},
		{
			name: "invalid name",
			fields: fields{
				user: &UsersServiceMock{
					CreateFunc: func(ctx context.Context, name string) (models.User, error) {
						return models.User{}, models.UserCreateParamInvalidNameErr
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{"name": ""}`)),
			},
			wantCode: http.StatusBadRequest,
			wantBody: []byte{0xa},
		},
		{
			name: "json decode failed",
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/", strings.NewReader(`{]`)),
			},
			wantCode: http.StatusInternalServerError,
			wantBody: []byte{0xa},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				user: tt.fields.user,
			}
			u.Create(tt.args.w, tt.args.r)
			assert.Equal(t, tt.wantCode, tt.args.w.Code, "invalid code")
			assert.Equal(t, tt.wantBody, tt.args.w.Body.Bytes(), "unxpected body")
		})
	}
}

func TestUsers_Fetch(t *testing.T) {
	type fields struct {
		user UsersService
	}
	type args struct {
		w *httptest.ResponseRecorder
		r *http.Request
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantCode int
		wantBody []byte
	}{
		{
			name: "success",
			fields: fields{
				user: &UsersServiceMock{
					FetchFunc: func(ctx context.Context, id string) (models.User, error) {
						// TODO: id should be taken from arg
						return models.User{ID: "1", Name: "mike"}, nil
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/1", strings.NewReader("")),
			},
			wantCode: http.StatusOK,
			wantBody: []byte(`{"id":"1","name":"mike"}` + "\n"),
		},
		{
			name: "not found",
			fields: fields{
				user: &UsersServiceMock{
					FetchFunc: func(ctx context.Context, id string) (models.User, error) {
						return models.User{}, models.NotFound
					},
				},
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/1", strings.NewReader("")),
			},
			wantCode: http.StatusNotFound,
			wantBody: []byte{0xa},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				user: tt.fields.user,
			}
			u.Fetch(tt.args.w, tt.args.r)
			assert.Equal(t, tt.wantCode, tt.args.w.Code, "invalid code")
			assert.Equal(t, tt.wantBody, tt.args.w.Body.Bytes(), "unxpected body")
		})
	}
}
