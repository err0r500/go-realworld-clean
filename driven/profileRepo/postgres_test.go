package profileRepo_test

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/err0r500/go-realworld-clean/driven/db"
	"github.com/err0r500/go-realworld-clean/driven/profileRepo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Start PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:18.0",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithStartupTimeout(60 * time.Second).
			WithOccurrence(2),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	require.NoError(t, err)

	host, err := container.Host(ctx)
	require.NoError(t, err)

	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	// Connect to database
	connStr := fmt.Sprintf(
		"host=%s port=%s user=testuser password=testpass dbname=testdb sslmode=disable",
		host, port.Port(),
	)

	database, err := sql.Open("pgx", connStr)
	require.NoError(t, err)

	// Wait for connection to be ready
	err = database.Ping()
	require.NoError(t, err)

	// Run schema creation directly (instead of migrations for simplicity in tests)
	schema := `
		CREATE TABLE IF NOT EXISTS follows (
			follower_id UUID NOT NULL,
			followee_id UUID NOT NULL,
			created_at TIMESTAMP DEFAULT NOW(),
			PRIMARY KEY (follower_id, followee_id),
			CHECK (follower_id != followee_id)
		);
		CREATE INDEX IF NOT EXISTS idx_follows_follower ON follows(follower_id);
		CREATE INDEX IF NOT EXISTS idx_follows_followee ON follows(followee_id);
	`
	_, err = database.Exec(schema)
	require.NoError(t, err)

	cleanup := func() {
		database.Close()
		container.Terminate(ctx)
	}

	return database, cleanup
}

// FIXME: these tests rely on their order
// (ex: idempotent rely on the previously created relationship)
func TestPostgresProfileRepo_Follow(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := profileRepo.NewPostgresProfileRepo(database)

	followerID := uuid.New().String()
	followeeID := uuid.New().String()

	t.Run("successfully create follow relationship", func(t *testing.T) {
		err := repo.Follow(followerID, followeeID)
		require.NoError(t, err)

		// Verify relationship exists
		isFollowing, err := repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.True(t, isFollowing)
	})

	t.Run("idempotent follow (no error on duplicate)", func(t *testing.T) {
		err := repo.Follow(followerID, followeeID)
		require.NoError(t, err)

		// Should still only have one relationship
		isFollowing, err := repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.True(t, isFollowing)
	})

	t.Run("cannot follow yourself", func(t *testing.T) {
		sameID := uuid.New().String()
		err := repo.Follow(sameID, sameID)
		require.Error(t, err)
		require.Contains(t, err.Error(), "cannot follow yourself")
	})
}

func TestPostgresProfileRepo_Unfollow(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := profileRepo.NewPostgresProfileRepo(database)

	followerID := uuid.New().String()
	followeeID := uuid.New().String()

	t.Run("successfully remove follow relationship", func(t *testing.T) {
		// Create relationship first
		err := repo.Follow(followerID, followeeID)
		require.NoError(t, err)

		// Verify it exists
		isFollowing, err := repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.True(t, isFollowing)

		// Unfollow
		err = repo.Unfollow(followerID, followeeID)
		require.NoError(t, err)

		// Verify it's gone
		isFollowing, err = repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.False(t, isFollowing)
	})

	t.Run("unfollow non-existent relationship (no error)", func(t *testing.T) {
		nonExistentFollower := uuid.New().String()
		nonExistentFollowee := uuid.New().String()

		err := repo.Unfollow(nonExistentFollower, nonExistentFollowee)
		require.NoError(t, err)
	})
}

func TestPostgresProfileRepo_IsFollowing(t *testing.T) {
	database, cleanup := setupTestDB(t)
	defer cleanup()

	repo := profileRepo.NewPostgresProfileRepo(database)

	followerID := uuid.New().String()
	followeeID := uuid.New().String()

	t.Run("returns false for non-existent relationship", func(t *testing.T) {
		isFollowing, err := repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.False(t, isFollowing)
	})

	t.Run("returns true for existing relationship", func(t *testing.T) {
		// Create relationship
		err := repo.Follow(followerID, followeeID)
		require.NoError(t, err)

		// Check it exists
		isFollowing, err := repo.IsFollowing(followerID, followeeID)
		require.NoError(t, err)
		require.True(t, isFollowing)
	})

	t.Run("follows are directional", func(t *testing.T) {
		user1 := uuid.New().String()
		user2 := uuid.New().String()

		// User1 follows User2
		err := repo.Follow(user1, user2)
		require.NoError(t, err)

		// User1 should follow User2
		isFollowing, err := repo.IsFollowing(user1, user2)
		require.NoError(t, err)
		require.True(t, isFollowing)

		// But User2 should NOT follow User1
		isFollowing, err = repo.IsFollowing(user2, user1)
		require.NoError(t, err)
		require.False(t, isFollowing)
	})
}

func init() {
	// Initialize db package to register pgx driver
	_, _ = db.NewPostgresConnection()
}
