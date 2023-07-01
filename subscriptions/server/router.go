package server

import "main/server/internal/handler"

func (s *Server) setupRoutes() {
	s.Router.Get("/_healthz", handler.Healthz)

	// http://localhost:8081/get?id=1
	s.Router.Get("/{id}", handler.UserSubscriptions(s.DB))
	s.Router.Get("/", handler.Subscriptions(s.DB))

	s.Router.Post("/subscribe", handler.Subscribe(s.DB))
	// TODO: route with url param without word "get" in url.
}
