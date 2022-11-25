package handlers

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
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
					CreateFunc: func(name string) (User, error) {
						if name != "mike" {
							t.Error("unexprected error")
						}
						return User{Name: name, ID: "1"}, nil
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := Users{
				user: tt.fields.user,
			}
			u.Create(tt.args.w, tt.args.r)
			if tt.args.w.Code != tt.wantCode {
				t.Errorf("unexpected code %d want %d", tt.args.w.Code, tt.wantCode)
			}
			if !reflect.DeepEqual(tt.args.w.Body.Bytes(), tt.wantBody) {
				t.Errorf("unexpected body %s, want %s", tt.args.w.Body.Bytes(), tt.wantBody)
			}
		})
	}
}
