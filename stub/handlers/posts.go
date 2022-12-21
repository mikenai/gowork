package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

type Post struct {
	ID        string
	Body      string
	Images    []string
	CreatedAt time.Time
}

func PostsHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(postsStub)
}

var postsStub = []Post{
	{
		ID: "123-1231-123123-123123",
		Images: []string{
			"https://cdn.stubserver.com/912451-124-124124-124124.jpg",
		},
		Body:      "Me and my wife heading to a pet store",
		CreatedAt: time.Now().Add(-time.Hour * 24),
	},
	{
		ID: "451-1231-12345-123123",
		Images: []string{
			"https://cdn.stubserver.com/451-1231-12345-123123.jpg",
			"https://cdn.stubserver.com/451-13681-12345-r5553.jpg",
		},
		Body:      "Today I planted new flowers. Aren't they awesome?",
		CreatedAt: time.Now().Add(-time.Hour * 24),
	},
}
