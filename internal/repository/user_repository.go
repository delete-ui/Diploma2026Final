package repository

import (
	"GolangBackendDiploma26/internal/models"
	"context"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *models.User) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (email, username, password_hash, role, is_verified, verification_code, verification_expires_at)
	          VALUES ($1, $2, $3, $4, $5, $6, $7)
	          RETURNING id, created_at, updated_at`

	err = tx.QueryRowContext(ctx, query,
		user.Email,
		user.Username,
		user.PasswordHash,
		user.Role,
		user.IsVerified,
		user.VerificationCode,
		user.VerificationExpiresAt,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert user: %w", err)
	}

	return tx.Commit()
}

func (r *UserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, role, is_verified,
	          verification_code, verification_expires_at, created_at, updated_at
	          FROM users WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user not found")
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, username, password_hash, role, is_verified,
	          verification_code, verification_expires_at, created_at, updated_at
	          FROM users WHERE email = $1`
	row := r.db.QueryRowContext(ctx, query, email)
	user := &models.User{}
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.VerificationCode,
		&user.VerificationExpiresAt,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("scan user: %w", err)
	}
	return user, nil
}

func (r *UserRepository) VerifyEmail(ctx context.Context, email, code string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET is_verified = TRUE, verification_code = NULL, verification_expires_at = NULL
	          WHERE email = $1 AND verification_code = $2 AND verification_expires_at > NOW() AND is_verified = FALSE`
	result, err := tx.ExecContext(ctx, query, email, code)
	if err != nil {
		return fmt.Errorf("update verification: %w", err)
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("invalid or expired verification code")
	}

	return tx.Commit()
}

func (r *UserRepository) SetResetCode(ctx context.Context, email, code string, expires time.Time) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET verification_code = $1, verification_expires_at = $2 WHERE email = $3`
	_, err = tx.ExecContext(ctx, query, code, expires, email)
	if err != nil {
		return fmt.Errorf("set reset code: %w", err)
	}
	return tx.Commit()
}

func (r *UserRepository) VerifyResetCode(ctx context.Context, email, code string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET verification_code = NULL, verification_expires_at = NULL
              WHERE email = $1 AND verification_code = $2 AND verification_expires_at > NOW()`
	result, err := tx.ExecContext(ctx, query, email, code)
	if err != nil {
		return fmt.Errorf("verify reset code: %w", err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("invalid or expired reset code")
	}
	return tx.Commit()
}

func (r *UserRepository) UpdatePassword(ctx context.Context, email, hashedPassword string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE users SET password_hash = $1, updated_at = NOW() WHERE email = $2`
	_, err = tx.ExecContext(ctx, query, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return tx.Commit()
}
