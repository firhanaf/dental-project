package service

import (
	"context"

	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
)

type BranchService struct{ repo *repository.BranchRepo }

func NewBranchService(repo *repository.BranchRepo) *BranchService {
	return &BranchService{repo: repo}
}

func (s *BranchService) List(ctx context.Context) ([]model.Branch, error) {
	return s.repo.List(ctx)
}
