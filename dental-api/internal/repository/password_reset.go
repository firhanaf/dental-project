package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type PasswordResetRepo struct{ db *pgxpool.Pool }

func NewPasswordResetRepo(db *pgxpool.Pool) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

func (r *PasswordResetRepo) InvalidatePrevious(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"DELETE FROM password_reset_tokens WHERE user_id=$1 AND used_at IS NULL",
		userID,
	)
	return err
}

func (r *PasswordResetRepo) Create(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.db.Exec(ctx,
		"INSERT INTO password_reset_tokens (user_id, token_hash, expires_at) VALUES ($1,$2,$3)",
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("create reset token: %w", err)
	}
	return nil
}

func (r *PasswordResetRepo) FindValid(ctx context.Context, userID uuid.UUID, tokenHash string) (*model.PasswordResetToken, error) {
	q := `SELECT id,user_id,token_hash,expires_at,used_at,created_at
	      FROM password_reset_tokens
	      WHERE user_id=$1 AND token_hash=$2 AND used_at IS NULL AND expires_at > NOW()`
	var t model.PasswordResetToken
	err := r.db.QueryRow(ctx, q, userID, tokenHash).Scan(
		&t.ID, &t.UserID, &t.TokenHash, &t.ExpiresAt, &t.UsedAt, &t.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("find valid token: %w", err)
	}
	return &t, nil
}

func (r *PasswordResetRepo) MarkUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE password_reset_tokens SET used_at=NOW() WHERE id=$1", id,
	)
	return err
}
