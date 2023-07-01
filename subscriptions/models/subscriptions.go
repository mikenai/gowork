package models

import (
	"fmt"
)

// Fetcher is an interface for getting a user's subscriptions.
type Fetcher interface {
	GetUserSubscriptions(userID string) ([]UserSubscription, error)
	GetSubscriptions() ([]UserSubscription, error)
}

// Subscriber is an interface for subscribing user.
type Subscriber interface {
	Subscribe(s UserSubscription) (string, error)
}

// UserSubscription contains information about a user's subscription.
type UserSubscription struct {
	UserID         string `json:"userId"`
	SubscriptionID string `json:"subscriptionId"`
	ChargeAmount   int    `json:"chargeAmount"`
	Status         string `json:"status"`
}

// Validate validates subscribe payload.
func Validate(s UserSubscription) error {

	if s.UserID == "" {
		return fmt.Errorf("UserID is missing")
	}

	if s.SubscriptionID == "" {
		return fmt.Errorf("SubscriptionID is missing")
	}

	if s.ChargeAmount == 0 {
		return fmt.Errorf("ChargeAmount is missing")
	}

	if s.Status == "" {
		return fmt.Errorf("Status is missing")
	}
	return nil
}
