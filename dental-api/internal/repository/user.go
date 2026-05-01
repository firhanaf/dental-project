package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type UserRepo struct{ db *pgxpool.Pool }
func NewUserRepo(db *pgxpool.Pool) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	q := `SELECT id,branch_id,name,email,password_hash,role,is_active,last_login_at,created_at,updated_at
	      FROM users WHERE email=$1 AND is_active=true`
	var u model.User
	err := r.db.QueryRow(ctx, q, email).Scan(
		&u.ID,&u.BranchID,&u.Name,&u.Email,&u.PasswordHash,
		&u.Role,&u.IsActive,&u.LastLoginAt,&u.CreatedAt,&u.UpdatedAt,
	)
	if err != nil { return nil, fmt.Errorf("find user: %w", err) }
	return &u, nil
}

func (r *UserRepo) FindByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	q := `SELECT id,branch_id,name,email,password_hash,role,is_active,last_login_at,created_at,updated_at
	      FROM users WHERE id=$1`
	var u model.User
	err := r.db.QueryRow(ctx, q, id).Scan(
		&u.ID,&u.BranchID,&u.Name,&u.Email,&u.PasswordHash,
		&u.Role,&u.IsActive,&u.LastLoginAt,&u.CreatedAt,&u.UpdatedAt,
	)
	if err != nil { return nil, fmt.Errorf("find user by id: %w", err) }
	return &u, nil
}

func (r *UserRepo) List(ctx context.Context) ([]model.User, error) {
	q := `SELECT id,branch_id,name,email,password_hash,role,is_active,last_login_at,created_at,updated_at
	      FROM users ORDER BY name ASC`
	rows, err := r.db.Query(ctx, q)
	if err != nil { return nil, err }
	defer rows.Close()
	users := make([]model.User, 0)
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID,&u.BranchID,&u.Name,&u.Email,&u.PasswordHash,
			&u.Role,&u.IsActive,&u.LastLoginAt,&u.CreatedAt,&u.UpdatedAt,
		); err != nil { return nil, err }
		users = append(users, u)
	}
	return users, nil
}

func (r *UserRepo) Create(ctx context.Context, u *model.User) error {
	u.ID = uuid.New()
	q := `INSERT INTO users (id,branch_id,name,email,password_hash,role)
	      VALUES ($1,$2,$3,$4,$5,$6) RETURNING created_at,updated_at`
	return r.db.QueryRow(ctx, q,
		u.ID,u.BranchID,u.Name,u.Email,u.PasswordHash,u.Role,
	).Scan(&u.CreatedAt,&u.UpdatedAt)
}

func (r *UserRepo) Update(ctx context.Context, u *model.User) error {
	q := `UPDATE users SET name=$1,role=$2,branch_id=$3,updated_at=NOW()
	      WHERE id=$4 RETURNING updated_at`
	return r.db.QueryRow(ctx, q, u.Name,u.Role,u.BranchID,u.ID).Scan(&u.UpdatedAt)
}

func (r *UserRepo) Deactivate(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE users SET is_active=false,updated_at=NOW() WHERE id=$1", id)
	return err
}

func (r *UserRepo) UpdateLastLogin(ctx context.Context, id uuid.UUID) error {
	now := time.Now()
	_, err := r.db.Exec(ctx,
		"UPDATE users SET last_login_at=$1 WHERE id=$2", now, id)
	return err
}
