package handler

import (
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
)

// Payload:
//	type UserSubscription struct {
//		UserID         string `json:"userId"`
//		SubscriptionID string `json:"subscriptionId"`
//		ChargeAmount   int    `json:"chargeAmount"`
//		Status         string `json:"status"`
//	}

func Subscribe(m models.Subscriber) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var payload models.UserSubscription
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("decode subscribe payload: %s, %w", r.RequestURI, err),
				http.StatusBadRequest,
				true,
			)
			return
		}

		err = models.Validate(payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error validating subscribe payload: %w", err),
				http.StatusBadRequest,
				true,
			)
			return
		}

		responseID, err := m.Subscribe(payload)
		if err != nil {
			handleError(
				w,
				fmt.Errorf("error subscribing user with user id: %s: %w", payload.UserID, err),
				http.StatusInternalServerError,
				true,
			)
			return
		}
		msg := fmt.Sprintf("User with ID: %s is signed up for subscription with ID: %s ", payload.UserID, responseID)
		response, _ := json.Marshal(msg)

		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
