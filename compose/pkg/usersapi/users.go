package usersapi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
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
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cl.BaseURL+fmt.Sprintf("/users/%s", id), nil)
	if err != nil {
		return User{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header["X-Request-ID"] = []string{middleware.GetReqID(ctx)}

	res, err := cl.Http.Do(req)
	if err != nil {
		return User{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return User{}, errors.New("non 200 code")
	}

	user := User{}
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return User{}, fmt.Errorf("json error: %w", err)
	}

	return user, nil
}
