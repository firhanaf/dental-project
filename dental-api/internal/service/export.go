package service

import (
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/repository"
	"github.com/yourusername/dental-api/pkg/excel"
)

type ExportService struct {
	patientRepo *repository.PatientRepo
	visitRepo   *repository.VisitRepo
}

func NewExportService(patientRepo *repository.PatientRepo, visitRepo *repository.VisitRepo) *ExportService {
	return &ExportService{patientRepo: patientRepo, visitRepo: visitRepo}
}

func (s *ExportService) ExportPatients(ctx context.Context, q dto.ExportQuery) (*bytes.Buffer, error) {
	var branchID *uuid.UUID
	if q.BranchID != "" {
		id, err := uuid.Parse(q.BranchID)
		if err == nil {
			branchID = &id
		}
	}

	patients, err := s.patientRepo.ListForExport(ctx, branchID)
	if err != nil {
		return nil, err
	}
	return excel.ExportPatients(patients)
}

func (s *ExportService) ExportVisits(ctx context.Context, q dto.ExportQuery) (*bytes.Buffer, error) {
	filter := repository.ExportVisitFilter{}

	if q.BranchID != "" {
		id, err := uuid.Parse(q.BranchID)
		if err == nil {
			filter.BranchID = &id
		}
	}
	if q.DateFrom != "" {
		t, err := time.Parse("2006-01-02", q.DateFrom)
		if err == nil {
			filter.DateFrom = &t
		}
	}
	if q.DateTo != "" {
		t, err := time.Parse("2006-01-02", q.DateTo)
		if err == nil {
			filter.DateTo = &t
		}
	}

	visits, err := s.visitRepo.ListForExport(ctx, filter)
	if err != nil {
		return nil, err
	}
	return excel.ExportVisits(visits)
}
