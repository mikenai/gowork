package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
)

type mockSubscriber struct {
	subscribe func(payload models.UserSubscription) (string, error)
}

func (m mockSubscriber) Subscribe(payload models.UserSubscription) (string, error) {
	return m.subscribe(payload)
}

func TestSubscribe(t *testing.T) {
	type args struct {
		subscriber models.Subscriber
		payload    models.UserSubscription
	}
	tests := []struct {
		name           string
		args           args
		expectedStatus int
	}{
		{
			name: "success",
			args: args{
				subscriber: mockSubscriber{
					subscribe: func(payload models.UserSubscription) (string, error) {
						return "1", nil
					},
				},
				payload: models.UserSubscription{
					UserID:         "1",
					SubscriptionID: "2",
					ChargeAmount:   3,
					Status:         "active",
				},
			},
			expectedStatus: 200,
		},
		{
			name: "bad request, no user id",
			args: args{
				subscriber: mockSubscriber{
					subscribe: func(payload models.UserSubscription) (string, error) {
						return "1", nil
					},
				},
				payload: models.UserSubscription{
					UserID:         "",
					SubscriptionID: "2",
					ChargeAmount:   3,
					Status:         "active",
				},
			},
			expectedStatus: 400,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			payload, err := json.Marshal(&tc.args.payload)
			if err != nil {
				t.Errorf("Error marshaling payload in test")
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/subscribe"), bytes.NewBuffer(payload))
			res := httptest.NewRecorder()

			router := chi.NewRouter()
			router.HandleFunc("/subscribe", Subscribe(tc.args.subscriber))

			router.ServeHTTP(res, req)

			if got := res.Code; got != tc.expectedStatus {
				t.Errorf("Unexpected response code: got %d, exp %d", got, tc.expectedStatus)
			}
		})
	}
}
