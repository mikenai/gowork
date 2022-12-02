package integrationtesting

import (
	"os"
	"testing"
)

func ShouldSkip(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") != "on" {
		t.Skip("integration test if off")
	}
}
