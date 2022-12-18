package users

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikenai/gowork/cmd/server/config"
	"github.com/mikenai/gowork/internal/models"
	integratiotesting "github.com/mikenai/gowork/pkg/integration_testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {

	m.Run()
}

func TestIntegrationStorage_Create(t *testing.T) {
	integratiotesting.ShouldSkip(t)

	db, dbClose := DbHelper(t)
	defer dbClose()

	SetUp(t, db, nil)
	t.Cleanup(func() {
		// clean after test
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	t.Run("succes", func(t *testing.T) {
		s := Storage{db: db}

		usr, err := s.Create(ctx, "mike")
		require.NoError(t, err)

		dbUser := models.User{}
		row := db.QueryRowContext(ctx, "SELECT id, name FROM users WHERE id=?", usr.ID)
		err = row.Scan(&dbUser.ID, &dbUser.Name)
		require.NoError(t, err)

		assert.Equal(t, usr, dbUser)
	})
}

func SetUp(t *testing.T, db *sql.DB, data any) {
	t.Helper()

	// place set up code
	// sql
}

func DbHelper(t *testing.T) (*sql.DB, func() error) {
	t.Helper()

	cfg, _, err := config.New()
	require.NoError(t, err)

	db, err := sql.Open("sqlite3", cfg.DB.DSN)
	require.NoError(t, err)

	return db, db.Close
}
