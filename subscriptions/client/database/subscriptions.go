package database

import "main/models"

// GetSubscriptions gets a list of all subscriptions from the database
func (c *Client) GetSubscriptions() (subscriptions []models.UserSubscription, err error) {
	rows, err := c.DB.Query("SELECT user_id, subscription_id, status, charge_amount FROM subscriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		subscription := models.UserSubscription{}
		if err := rows.Scan(&subscription.UserID, &subscription.SubscriptionID, &subscription.Status, &subscription.ChargeAmount); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

// GetUserSubscriptions gets a list of a user's subscriptions from the database
func (c *Client) GetUserSubscriptions(userID string) (subscriptions []models.UserSubscription, err error) {
	sqlStatement := "SELECT user_id, subscription_id, status, charge_amount FROM subscriptions WHERE user_id = $1"

	rows, err := c.DB.Query(sqlStatement, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		subscription := models.UserSubscription{}
		if err := rows.Scan(&subscription.UserID, &subscription.SubscriptionID, &subscription.Status, &subscription.ChargeAmount); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, err
}
