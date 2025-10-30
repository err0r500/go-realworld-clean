package profileRepo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

// PostgresProfileRepo implements the ProfileRepository interface using PostgreSQL
type PostgresProfileRepo struct {
	db *sql.DB
}

// NewPostgresProfileRepo creates a new PostgreSQL profile repository
func NewPostgresProfileRepo(db *sql.DB) *PostgresProfileRepo {
	return &PostgresProfileRepo{db: db}
}

// Follow creates a follow relationship between follower and followee
func (r *PostgresProfileRepo) Follow(followerID, followeeID string) error {
	if followerID == followeeID {
		return fmt.Errorf("cannot follow yourself")
	}

	query := `
		INSERT INTO follows (follower_id, followee_id)
		VALUES ($1, $2)
		ON CONFLICT (follower_id, followee_id) DO NOTHING
	`

	_, err := r.db.Exec(query, followerID, followeeID)
	if err != nil {
		// Check for constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23514" { // check_violation
				return fmt.Errorf("cannot follow yourself")
			}
		}
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	return nil
}

// Unfollow removes a follow relationship between follower and followee
func (r *PostgresProfileRepo) Unfollow(followerID, followeeID string) error {
	query := `
		DELETE FROM follows
		WHERE follower_id = $1 AND followee_id = $2
	`

	result, err := r.db.Exec(query, followerID, followeeID)
	if err != nil {
		return fmt.Errorf("failed to remove follow relationship: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		// Not an error, just means relationship didn't exist
		return nil
	}

	return nil
}

// IsFollowing checks if follower is following followee
func (r *PostgresProfileRepo) IsFollowing(followerID, followeeID string) (bool, error) {
	query := `
		SELECT EXISTS(
			SELECT 1 FROM follows
			WHERE follower_id = $1 AND followee_id = $2
		)`

	var exists bool
	err := r.db.QueryRow(query, followerID, followeeID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check follow relationship: %w", err)
	}

	return exists, nil
}
