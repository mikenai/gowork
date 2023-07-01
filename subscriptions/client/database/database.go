// Package database contains a Postgres client and methods for communicating with the database.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"main/config"
)

// https://medium.com/gojekengineering/adding-a-database-to-a-go-web-application-b0e8e8b16fb9

// Client holds the database client and prepared statements.
type Client struct {
	DB *sql.DB
}

const (
	DatabaseUrl      = "localhost"
	DatabasePort     = 5432
	DatabaseUser     = "postgres"
	DatabasePassword = "postgres"
	DatabaseDb       = "subscriptions"
)

// Init sets up a new database client.
func (c *Client) Init(ctx context.Context, config *config.Config) error {
	connString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		DatabaseUrl, DatabasePort, DatabaseUser, DatabasePassword, DatabaseDb)

	db, err := sql.Open("postgres", connString)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	c.DB = db
	return nil
}

// Close closes the database connection and statements.
func (c *Client) Close() error {
	if err := c.DB.Close(); err != nil {
		return err
	}
	return nil
}

const active = "ACTIVE"
