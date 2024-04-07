package main

import (
	"github.com/stretchr/testify/assert"
	cfg "github.com/tangvis/erp/conf/config"
	"os"
	"testing"
)

// TestInitDB tests the initDB function
func TestInitDB(t *testing.T) {
	_ = os.Setenv(cfg.EnvKey, "dev")
	// Setup
	initConfig()

	// Act
	db := initDB() // Assuming initDB doesn't return an error but panics on failure

	// Assert
	// Here you'd assert that db is correctly initialized based on your mock.
	// This is more about conceptual understanding since you can't easily
	// check internal state without an actual DB connection or further mocking.

	assert.NotNil(t, db, "DB should not be nil")
}
