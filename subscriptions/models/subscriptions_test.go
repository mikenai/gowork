package models

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Validate(t *testing.T) {
	tests := []struct {
		name       string
		dataStruct UserSubscription
		err        error
	}{
		{
			name:       "success",
			dataStruct: UserSubscription{UserID: "1", SubscriptionID: "2", ChargeAmount: 3, Status: "active"},
			err:        nil,
		},
		{
			name:       "error: empty charge amount",
			dataStruct: UserSubscription{UserID: "1", SubscriptionID: "2", ChargeAmount: 0, Status: "active"},
			err:        fmt.Errorf("ChargeAmount is missing"),
		},
		{
			name:       "error: empty user id",
			dataStruct: UserSubscription{UserID: "", SubscriptionID: "2", ChargeAmount: 1, Status: "active"},
			err:        fmt.Errorf("UserID is missing"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.dataStruct)
			assert.Equal(t, tt.err, err)
		})
	}
}
