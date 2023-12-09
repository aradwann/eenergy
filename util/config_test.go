package util

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfigFromFile(t *testing.T) {
	// Set up a temporary configuration file for testing
	configContent := `
	ENVIRONMENT=test
	DB_DRIVER=postgres
	DB_SOURCE=example-db-source
	MIGRATIONS_URL=example-migrations-url
	HTTP_SERVER_ADDRESS=localhost:8080
	REDIS_ADDRESS=localhost:6379
	TOKEN_SYMMETRIC_KEY=example-key
	ACCESS_TOKEN_DURATION=1h
	REFRESH_TOKEN_DURATION=24h
	`
	tempFile, err := createTempConfigFile(configContent)
	assert.NoError(t, err)
	defer os.Remove(tempFile.Name())

	// Load configuration from the temporary file
	config, err := LoadConfig(".", tempFile.Name())
	assert.NoError(t, err)

	// Assert the expected values
	assert.Equal(t, "test", config.Environment)
	assert.Equal(t, "postgres", config.DBDriver)
	assert.Equal(t, "example-db-source", config.DBSource)
	assert.Equal(t, "example-migrations-url", config.MigrationsURL)
	assert.Equal(t, "localhost:8080", config.HTTPServerAddress)
	assert.Equal(t, "localhost:6379", config.RedisAddress)
	assert.Equal(t, "example-key", config.TokenSymmetricKey)
	assert.Equal(t, 1*time.Hour, config.AccessTokenDuration)
	assert.Equal(t, 24*time.Hour, config.RefreshTokenDuration)
}

func createTempConfigFile(content string) (*os.File, error) {
	tempFile, err := os.CreateTemp(".", "app.env")
	if err != nil {
		return nil, err
	}

	_, err = tempFile.WriteString(content)
	if err != nil {
		tempFile.Close()
		return nil, err
	}

	err = tempFile.Close()
	if err != nil {
		return nil, err
	}

	return tempFile, nil
}
