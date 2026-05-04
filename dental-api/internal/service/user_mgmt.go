package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
)

type UserMgmtService struct{ repo *repository.UserRepo }

func NewUserMgmtService(repo *repository.UserRepo) *UserMgmtService {
	return &UserMgmtService{repo: repo}
}

func (s *UserMgmtService) List(ctx context.Context) ([]model.User, error) {
	return s.repo.List(ctx)
}

func (s *UserMgmtService) Create(ctx context.Context, req *dto.CreateUserRequest) (*model.User, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama wajib diisi")
	}
	if req.Email == "" {
		return nil, fmt.Errorf("email wajib diisi")
	}
	if len(req.Password) < 8 {
		return nil, fmt.Errorf("password minimal 8 karakter")
	}
	role := model.UserRole(req.Role)
	if role != model.RoleSuperAdmin && role != model.RoleWrite && role != model.RoleReadonly {
		return nil, fmt.Errorf("role tidak valid, gunakan 'superadmin', 'write', atau 'readonly'")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	u := &model.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	if req.BranchID != nil && *req.BranchID != "" {
		id, err := uuid.Parse(*req.BranchID)
		if err != nil {
			return nil, fmt.Errorf("branch_id tidak valid")
		}
		u.BranchID = &id
	}

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return u, nil
}

func (s *UserMgmtService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateUserRequest) (*model.User, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama wajib diisi")
	}
	role := model.UserRole(req.Role)
	if role != model.RoleSuperAdmin && role != model.RoleWrite && role != model.RoleReadonly {
		return nil, fmt.Errorf("role tidak valid")
	}

	u := &model.User{ID: id, Name: req.Name, Role: role}

	if req.BranchID != nil && *req.BranchID != "" {
		branchID, err := uuid.Parse(*req.BranchID)
		if err != nil {
			return nil, fmt.Errorf("branch_id tidak valid")
		}
		u.BranchID = &branchID
	}

	if err := s.repo.Update(ctx, u); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update user: %w", err)
	}
	return u, nil
}

func (s *UserMgmtService) Deactivate(ctx context.Context, id uuid.UUID) error {
	return s.repo.Deactivate(ctx, id)
}

func (s *UserMgmtService) Activate(ctx context.Context, id uuid.UUID) error {
	return s.repo.Activate(ctx, id)
}

func (s *UserMgmtService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return fmt.Errorf("user tidak dapat dihapus karena memiliki riwayat kunjungan — gunakan Nonaktifkan")
		}
		return err
	}
	return nil
}
