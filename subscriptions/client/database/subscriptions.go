package database

import "main/models"

func (c *Client) GetSubscriptions() (subscriptions []models.UserSubscription, err error) {

	rows, err := c.DB.Query("SELECT subscription_id, status, charge_amount FROM subscriptions")
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
	return subscriptions, nil
}

func (c *Client) Subscribe(s models.UserSubscription) (subscriptionID string, err error) {
	sqlStatementINSERT := "INSERT INTO subscriptions(user_id, subscription_id, status, charge_amount) VALUES ($1,$2,$3,$4) RETURNING subscription_id"

	err = c.DB.QueryRow(sqlStatementINSERT, s.UserID, s.SubscriptionID, s.Status, s.ChargeAmount).Scan(&s.SubscriptionID)
	if err != nil {
		return subscriptionID, err
	}
	return subscriptionID, nil
}
