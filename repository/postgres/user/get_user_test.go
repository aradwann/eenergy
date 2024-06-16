package user

import (
	"context"
	"database/sql"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/aradwann/eenergy/repository/postgres/migrate"
	"github.com/aradwann/eenergy/telemetry"
	"github.com/aradwann/eenergy/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var db *sql.DB

func setupTestContainer() (func(), *sql.DB, error) {
	ctx := context.Background()

	dbName := "eenergy"
	dbUser := "user"
	dbPassword := "password"

	postgresC, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgis/postgis:16-3.4-alpine"),
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return nil, nil, err
	}

	connStr, err := postgresC.ConnectionString(ctx)
	if err != nil {
		return nil, nil, err
	}
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, nil, err
	}
	// Load configuration.
	_, err = util.LoadConfig("../../..", "test.env")
	if err != nil {
		return nil, nil, err
	}
	// Ensure the migrations URL is correct and use an absolute path
	migrationsPath, err := filepath.Abs("/home/aradwan/Documents/repos/eenergy/migrations")
	if err != nil {
		slog.Error("Failed to determine absolute path", slog.String("error", err.Error()))
		os.Exit(1)
	}
	migrationsURL := "file://" + migrationsPath
	// Run any migrations on the database
	migrate.RunDBMigrations(db, migrationsURL)

	// 2. Create a snapshot of the database to restore later
	err = postgresC.Snapshot(ctx, postgres.WithSnapshotName("test-snapshot"))
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		db.Close()
		postgresC.Terminate(ctx)
	}

	return cleanup, db, nil
}

func TestMain(m *testing.M) {
	var err error
	var cleanup func()

	cleanup, db, err = setupTestContainer()
	if err != nil {
		log.Fatalf("Could not set up test container: %s", err)
	}

	code := m.Run()

	cleanup()
	os.Exit(code)
}

func TestGetUser(t *testing.T) {
	config := util.Config{
		Environment: "test",
	}
	// Setup logger (assuming slog.Logger is already initialized properly)
	logger := telemetry.InitLogger(config)

	repo := NewUserRepository(db, logger)

	// Prepare the database for the test
	ctx := context.Background()

	// Test GetUser method
	user, err := repo.GetUser(ctx, "testuser")
	require.NoError(t, err)

	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "testuser@example.com", user.Email)
}
