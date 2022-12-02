package users

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mikenai/gowork/internal/models"
	integratiotesting "github.com/mikenai/gowork/pkg/integration_testing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntegrationStorage_Create(t *testing.T) {
	integratiotesting.ShouldSkip(t)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := sql.Open("sqlite3", "/Users/mikenai/dev/code/gowork/tmp/db.sqlite3")
	require.NoError(t, err)

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
