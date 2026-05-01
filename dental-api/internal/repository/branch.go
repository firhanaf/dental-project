package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type BranchRepo struct{ db *pgxpool.Pool }

func NewBranchRepo(db *pgxpool.Pool) *BranchRepo { return &BranchRepo{db: db} }

func (r *BranchRepo) List(ctx context.Context) ([]model.Branch, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, name, address, phone, is_active, created_at, updated_at
		 FROM branches ORDER BY name`)
	if err != nil {
		return nil, fmt.Errorf("list branches: %w", err)
	}
	defer rows.Close()

	branches := make([]model.Branch, 0)
	for rows.Next() {
		var b model.Branch
		if err := rows.Scan(&b.ID, &b.Name, &b.Address, &b.Phone, &b.IsActive, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan branch: %w", err)
		}
		branches = append(branches, b)
	}
	return branches, nil
}
