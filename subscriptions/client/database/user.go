package database

import (
	"main/models"
)

// GetUserSubscriptions gets a list of a user's subscriptions from the database,
// returning all active subscriptions as well as the most recently cancelled
// subscription for each product.
func (c *Client) GetUserSubscriptions(userID string) (subscriptions []models.UserSubscription, err error) {
	sqlStatement := "SELECT subscription_id, status, charge_amount FROM subscriptions WHERE user_id = $1"

	rows, err := c.DB.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		subscription := models.UserSubscription{}
		if err := rows.Scan(&subscription.SubscriptionID, &subscription.Status, &subscription.ChargeAmount); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, err
}
