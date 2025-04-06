package db

import (
	"testing"

	"github.com/DebroyeAntoine/go_link_vault/internal/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConnect(t *testing.T) {
	// Test normal connection
	SetupTestDB()

	// Test if DB is connected
	if DB == nil {
		t.Fatal("DB connection failed")
	}

	// Test AutoMigrate
	err := DB.AutoMigrate(&models.Link{})
	assert.NoError(t, err, "Migration should succeed")

	// Test if the 'links' table exists in the DB
	var links []models.Link
	result := DB.Find(&links)
	assert.NoError(t, result.Error, "Querying links table should succeed")
	assert.Equal(t, len(links), 0, "Newly created links table should be empty")
}

func TestConnectFail(t *testing.T) {
	// Intentionally breaking the connection to test failure
	dsn := "host=invalid_host user=invalid_user password=invalid dbname=invalid port=5432 sslmode=disable"
	_, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	assert.Error(t, err, "Should fail to connect to DB with incorrect credentials")
}
