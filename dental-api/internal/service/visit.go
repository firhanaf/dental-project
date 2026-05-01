package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
)

type VisitService struct {
	repo        *repository.VisitRepo
	patientRepo *repository.PatientRepo
}

func NewVisitService(repo *repository.VisitRepo, patientRepo *repository.PatientRepo) *VisitService {
	return &VisitService{repo: repo, patientRepo: patientRepo}
}

func (s *VisitService) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]model.Visit, error) {
	return s.repo.ListByPatient(ctx, patientID)
}

func (s *VisitService) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	v, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return v, nil
}

func (s *VisitService) Create(ctx context.Context, req *dto.CreateVisitRequest, claims *jwtpkg.Claims) (*model.Visit, error) {
	if req.PatientID == "" {
		return nil, fmt.Errorf("patient_id wajib diisi")
	}
	if req.DoctorID == "" {
		return nil, fmt.Errorf("doctor_id wajib diisi")
	}
	if req.VisitDate == "" {
		return nil, fmt.Errorf("visit_date wajib diisi")
	}
	if req.ChiefComplaint == "" {
		return nil, fmt.Errorf("chief_complaint wajib diisi")
	}

	patientID, err := uuid.Parse(req.PatientID)
	if err != nil {
		return nil, fmt.Errorf("patient_id tidak valid")
	}
	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("doctor_id tidak valid")
	}
	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return nil, fmt.Errorf("format visit_date tidak valid, gunakan YYYY-MM-DD")
	}

	// branch_id: dari claims (write) atau dari data pasien (superadmin)
	var branchID uuid.UUID
	if claims.BranchID != nil {
		branchID = *claims.BranchID
	} else {
		patient, err := s.patientRepo.GetByID(ctx, patientID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("pasien tidak ditemukan")
			}
			return nil, fmt.Errorf("cari pasien: %w", err)
		}
		branchID = patient.BranchID
	}

	var nextControl *time.Time
	if req.NextControlDate != nil && *req.NextControlDate != "" {
		t, err := time.Parse("2006-01-02", *req.NextControlDate)
		if err != nil {
			return nil, fmt.Errorf("format next_control_date tidak valid, gunakan YYYY-MM-DD")
		}
		nextControl = &t
	}

	v := &model.Visit{
		PatientID:       patientID,
		BranchID:        branchID,
		DoctorID:        doctorID,
		CreatedBy:       claims.UserID,
		VisitDate:       visitDate,
		ChiefComplaint:  req.ChiefComplaint,
		Diagnosis:       req.Diagnosis,
		Treatment:       req.Treatment,
		TeethInvolved:   req.TeethInvolved,
		Cost:            req.Cost,
		NextControlDate: nextControl,
		Notes:           req.Notes,
	}

	if err := s.repo.Create(ctx, v); err != nil {
		return nil, fmt.Errorf("create visit: %w", err)
	}
	return v, nil
}

func (s *VisitService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdateVisitRequest) (*model.Visit, error) {
	if req.DoctorID == "" {
		return nil, fmt.Errorf("doctor_id wajib diisi")
	}
	if req.VisitDate == "" {
		return nil, fmt.Errorf("visit_date wajib diisi")
	}
	if req.ChiefComplaint == "" {
		return nil, fmt.Errorf("chief_complaint wajib diisi")
	}

	doctorID, err := uuid.Parse(req.DoctorID)
	if err != nil {
		return nil, fmt.Errorf("doctor_id tidak valid")
	}
	visitDate, err := time.Parse("2006-01-02", req.VisitDate)
	if err != nil {
		return nil, fmt.Errorf("format visit_date tidak valid, gunakan YYYY-MM-DD")
	}

	var nextControl *time.Time
	if req.NextControlDate != nil && *req.NextControlDate != "" {
		t, err := time.Parse("2006-01-02", *req.NextControlDate)
		if err != nil {
			return nil, fmt.Errorf("format next_control_date tidak valid")
		}
		nextControl = &t
	}

	v := &model.Visit{
		ID:              id,
		DoctorID:        doctorID,
		VisitDate:       visitDate,
		ChiefComplaint:  req.ChiefComplaint,
		Diagnosis:       req.Diagnosis,
		Treatment:       req.Treatment,
		TeethInvolved:   req.TeethInvolved,
		Cost:            req.Cost,
		NextControlDate: nextControl,
		Notes:           req.Notes,
	}

	if err := s.repo.Update(ctx, v); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update visit: %w", err)
	}
	return v, nil
}

func (s *VisitService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.repo.SoftDelete(ctx, id)
	if err != nil && err.Error() == "record not found" {
		return ErrNotFound
	}
	return err
}
