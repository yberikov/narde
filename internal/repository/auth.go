package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gofrs/uuid"
	"narde/internal/domain"

	"github.com/jmoiron/sqlx" // <-- Import sqlx
)

type AuthRepository struct {
	db *sqlx.DB // <-- Use sqlx.DB
}

// NewAuthPostgresRepo creates a new auth repository.
func NewAuthRepository(db *sqlx.DB) *AuthRepository { // <-- Accept sqlx.DB
	return &AuthRepository{db: db}
}

func (r *AuthRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at`

	// Use QueryRowxContext for consistency
	err := r.db.QueryRowxContext(ctx, query, user.Username, user.Email, user.Password).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUserByUsername finds a user by their username.
func (r *AuthRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE username = $1`

	// Use GetContext to scan the whole struct at once
	err := r.db.GetContext(ctx, &user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound // A custom error is better here
		}
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	return &user, nil
}

// GetUserByID finds a user by their unique ID.
func (r *AuthRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*domain.User, error) {
	var user domain.User
	query := `SELECT * FROM users WHERE id = $1`

	// Use GetContext here as well
	err := r.db.GetContext(ctx, &user, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}
