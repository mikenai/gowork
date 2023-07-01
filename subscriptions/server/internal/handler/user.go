package handler

import (
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
)

func UserSubscriptions(m models.Fetcher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID
		userID := r.URL.Query().Get("id")
		if userID == "" {
			handleError(
				w,
				fmt.Errorf("missing user ID"),
				http.StatusBadRequest,
				true)
			return
		}
		subscriptions, err := m.GetUserSubscriptions(userID)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error getting user subscriptions: %w", err),
				http.StatusInternalServerError,
				true,
			)
			return
		}

		response, _ := json.Marshal(subscriptions)

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
