package stub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type Post struct {
	ID        string
	Body      string
	Images    []string
	CreatedAt time.Time
}

type Profile struct {
	Picture string
	Bio     string
}

type Client struct {
	BaseURL string
	Http    http.Client
}

func (cl *Client) GetPosts(ctx context.Context, id string) ([]Post, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		cl.BaseURL+fmt.Sprintf("/%s/posts", id), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header["X-Request-ID"] = []string{middleware.GetReqID(ctx)}

	res, err := cl.Http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("non 200 code")
	}

	var user []Post
	if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("json error: %w", err)
	}

	return user, nil
}

func (cl *Client) GetProfile(ctx context.Context, id string) (Profile, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, cl.BaseURL+fmt.Sprintf("/profiles/%s", id), nil)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header["X-Request-ID"] = []string{middleware.GetReqID(ctx)}

	res, err := cl.Http.Do(req)
	if err != nil {
		return Profile{}, fmt.Errorf("failed to perform request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Profile{}, errors.New("non 200 code")
	}

	profile := Profile{}
	if err := json.NewDecoder(res.Body).Decode(&profile); err != nil {
		return Profile{}, fmt.Errorf("json error: %w", err)
	}

	return profile, nil
}
