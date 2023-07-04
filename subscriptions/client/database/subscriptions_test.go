package database

import (
	"database/sql"
	"main/models"
	"testing"

	_ "github.com/lib/pq"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptions(t *testing.T) {
	db, dbClose := DbHelper(t)
	defer dbClose()

	SetUp(t, db, nil)
	t.Cleanup(func() {
	})

	t.Run("succes", func(t *testing.T) {
		s := Client{DB: db}

		subscriptions, err := s.GetSubscriptions()
		require.NoError(t, err)

		dbSubscriptions := []models.UserSubscription{}
		rows, err := db.Query("SELECT user_id, subscription_id, status, charge_amount FROM subscriptions")
		require.NoError(t, err)

		for rows.Next() {
			dbSubscription := models.UserSubscription{}
			err := rows.Scan(&dbSubscription.UserID, &dbSubscription.SubscriptionID, &dbSubscription.Status, &dbSubscription.ChargeAmount)
			require.NoError(t, err)
			dbSubscriptions = append(dbSubscriptions, dbSubscription)
		}
		assert.Equal(t, subscriptions, dbSubscriptions)
	})

}

func TestGetUserSubscriptions(t *testing.T) {
	db, dbClose := DbHelper(t)
	defer dbClose()

	SetUp(t, db, nil)
	t.Cleanup(func() {
	})

	t.Run("succes", func(t *testing.T) {
		s := Client{DB: db}

		subscriptions, err := s.GetUserSubscriptions("1")
		require.NoError(t, err)

		dbSubscriptions := []models.UserSubscription{}
		sqlStatement := "SELECT user_id, subscription_id, status, charge_amount FROM subscriptions WHERE user_id = $1"

		rows, err := db.Query(sqlStatement, "1")
		require.NoError(t, err)
		defer rows.Close()

		for rows.Next() {
			dbSubscription := models.UserSubscription{}
			err := rows.Scan(&dbSubscription.UserID, &dbSubscription.SubscriptionID, &dbSubscription.Status, &dbSubscription.ChargeAmount)
			require.NoError(t, err)
			dbSubscriptions = append(dbSubscriptions, dbSubscription)
		}
		assert.Equal(t, subscriptions, dbSubscriptions)
	})

}

func SetUp(t *testing.T, db *sql.DB, data any) {
	t.Helper()
}

func DbHelper(t *testing.T) (*sql.DB, func() error) {
	t.Helper()
	db, err := sql.Open("postgres", ConnString())
	require.NoError(t, err)

	return db, db.Close
}
