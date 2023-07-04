package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
)

type mockFetcher struct {
	getUserSubscriptions func(userID string) ([]models.UserSubscription, error)
	getSubscriptions     func() ([]models.UserSubscription, error)
}

func (m mockFetcher) GetUserSubscriptions(userID string) ([]models.UserSubscription, error) {
	return m.getUserSubscriptions(userID)
}

func (m mockFetcher) GetSubscriptions() ([]models.UserSubscription, error) {
	return m.getSubscriptions()
}

func TestGetUserSubscriptions(t *testing.T) {
	type args struct {
		fetcher models.Fetcher
		userID  string
	}
	test := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "success",
			args: args{
				fetcher: mockFetcher{
					getUserSubscriptions: func(userID string) (subscriptions []models.UserSubscription, err error) {
						assert.Equal(t, "1", userID)
						return subscriptions, nil
					},
				},
				userID: "1",
			},
			expectedStatus: 200,
			expectedBody:   []byte(`[{"userId":"1","subscriptionId":"5678","chargeAmount":45,"status":"1"}]`),
		},
		{
			name: "no user id",
			args: args{
				fetcher: mockFetcher{
					getUserSubscriptions: func(userID string) (subscriptions []models.UserSubscription, err error) {
						assert.Equal(t, "1", userID)
						return subscriptions, nil
					},
				},
				userID: "",
			},
			expectedStatus: 400,
			expectedBody:   nil,
		},
		{
			name: "error getting user subscriptions",
			args: args{
				fetcher: mockFetcher{
					getUserSubscriptions: func(userID string) (subscriptions []models.UserSubscription, err error) {
						return nil, errors.New("error")
					},
				},
				userID: "2",
			},
			expectedStatus: 500,
			expectedBody:   nil,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/get?id=%s", tc.args.userID), nil)
			res := httptest.NewRecorder()

			router := chi.NewRouter()
			router.HandleFunc("/{id}", UserSubscriptions(tc.args.fetcher))

			router.ServeHTTP(res, req)
			var s models.UserSubscription
			var o models.UserSubscription
			_ = json.Unmarshal(res.Body.Bytes(), &s)

			_ = json.Unmarshal(tc.expectedBody, &o)

			assert.Equal(t, tc.expectedStatus, res.Code)
			assert.Equal(t, s, o, "unxpected body")
		})
	}
}

func TestGetSubscriptions(t *testing.T) {
	type args struct {
		fetcher models.Fetcher
	}
	test := []struct {
		name           string
		args           args
		expectedStatus int
		expectedBody   []byte
	}{
		{
			name: "success",
			args: args{
				fetcher: mockFetcher{
					getSubscriptions: func() (subscriptions []models.UserSubscription, err error) {
						return subscriptions, nil
					},
				},
			},
			expectedStatus: 200,
			expectedBody:   []byte(`[{"userId":"1","subscriptionId":"1234","chargeAmount":124,"status":"1"},{"userId":"2","subscriptionId":"5678","chargeAmount":45,"status":"1"}]`),
		},
		{
			name: "error getting subscriptions",
			args: args{
				fetcher: mockFetcher{
					getSubscriptions: func() (subscriptions []models.UserSubscription, err error) {
						return nil, errors.New("error")
					},
				},
			},
			expectedStatus: 500,
			expectedBody:   nil,
		},
	}

	for _, tc := range test {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/"), nil)
			res := httptest.NewRecorder()

			router := chi.NewRouter()
			router.HandleFunc("/", Subscriptions(tc.args.fetcher))

			router.ServeHTTP(res, req)
			var s models.UserSubscription
			var o models.UserSubscription
			_ = json.Unmarshal(res.Body.Bytes(), &s)

			_ = json.Unmarshal(tc.expectedBody, &o)

			assert.Equal(t, tc.expectedStatus, res.Code)
			assert.Equal(t, s, o, "unxpected body")
		})
	}
}
