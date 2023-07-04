package database

import (
	"main/models"
)

// Subscribe creates inserts user subscription information into the database.
func (c *Client) Subscribe(s models.UserSubscription) (subscriptionID string, err error) {
	sqlStatementINSERT := "INSERT INTO subscriptions(user_id, subscription_id, status, charge_amount) VALUES ($1,$2,$3,$4) RETURNING subscription_id"

	err = c.DB.QueryRow(sqlStatementINSERT, s.UserID, s.SubscriptionID, s.Status, s.ChargeAmount).Scan(&s.SubscriptionID)
	if err != nil {
		return subscriptionID, err
	}
	return subscriptionID, nil
}
