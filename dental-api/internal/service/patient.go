package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"unicode"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
)

func validateNIK(nik string) error {
	if len(nik) != 16 {
		return fmt.Errorf("NIK harus 16 digit angka")
	}
	for _, c := range nik {
		if !unicode.IsDigit(c) {
			return fmt.Errorf("NIK harus 16 digit angka")
		}
	}
	return nil
}

type PatientService struct{ repo *repository.PatientRepo }

func NewPatientService(repo *repository.PatientRepo) *PatientService {
	return &PatientService{repo: repo}
}

func (s *PatientService) List(ctx context.Context, f repository.PatientFilter) ([]model.PatientListRow, int, error) {
	return s.repo.List(ctx, f)
}

func (s *PatientService) GetByID(ctx context.Context, id uuid.UUID) (*model.Patient, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return p, nil
}

func (s *PatientService) Create(ctx context.Context, req *dto.CreatePatientRequest, claims *jwtpkg.Claims) (*model.Patient, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama wajib diisi")
	}
	if req.DateOfBirth == "" {
		return nil, fmt.Errorf("tanggal lahir wajib diisi")
	}
	if req.Gender != string(model.GenderMale) && req.Gender != string(model.GenderFemale) {
		return nil, fmt.Errorf("gender tidak valid, gunakan 'male' atau 'female'")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("nomor telepon wajib diisi")
	}

	var branchID uuid.UUID
	if claims.Role == model.RoleSuperAdmin {
		if req.BranchID == "" {
			return nil, fmt.Errorf("branch_id wajib diisi")
		}
		id, err := uuid.Parse(req.BranchID)
		if err != nil {
			return nil, fmt.Errorf("branch_id tidak valid")
		}
		branchID = id
	} else {
		if claims.BranchID == nil {
			return nil, fmt.Errorf("branch tidak ditemukan pada akun ini")
		}
		branchID = *claims.BranchID
	}

	// Validasi NIK jika diisi
	if req.NIK != nil && *req.NIK != "" {
		if err := validateNIK(*req.NIK); err != nil {
			return nil, err
		}
		exists, err := s.repo.ExistsByNIK(ctx, *req.NIK, nil)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("NIK sudah terdaftar untuk pasien lain")
		}
	} else {
		req.NIK = nil // normalisasi string kosong menjadi nil
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
	}

	p := &model.Patient{
		BranchID:     branchID,
		CreatedBy:    claims.UserID,
		Name:         req.Name,
		NIK:          req.NIK,
		DateOfBirth:  dob,
		Gender:       model.GenderType(req.Gender),
		Phone:        req.Phone,
		Address:      req.Address,
		Occupation:   req.Occupation,
		AllergyNotes: req.AllergyNotes,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("create patient: %w", err)
	}
	return p, nil
}

func (s *PatientService) Update(ctx context.Context, id uuid.UUID, req *dto.UpdatePatientRequest) (*model.Patient, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("nama wajib diisi")
	}
	if req.Gender != string(model.GenderMale) && req.Gender != string(model.GenderFemale) {
		return nil, fmt.Errorf("gender tidak valid")
	}
	if req.Phone == "" {
		return nil, fmt.Errorf("nomor telepon wajib diisi")
	}

	// Validasi NIK jika diisi
	if req.NIK != nil && *req.NIK != "" {
		if err := validateNIK(*req.NIK); err != nil {
			return nil, err
		}
		exists, err := s.repo.ExistsByNIK(ctx, *req.NIK, &id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, fmt.Errorf("NIK sudah terdaftar untuk pasien lain")
		}
	} else {
		req.NIK = nil
	}

	dob, err := time.Parse("2006-01-02", req.DateOfBirth)
	if err != nil {
		return nil, fmt.Errorf("format tanggal lahir tidak valid, gunakan YYYY-MM-DD")
	}

	p := &model.Patient{
		ID:           id,
		Name:         req.Name,
		NIK:          req.NIK,
		DateOfBirth:  dob,
		Gender:       model.GenderType(req.Gender),
		Phone:        req.Phone,
		Address:      req.Address,
		Occupation:   req.Occupation,
		AllergyNotes: req.AllergyNotes,
	}

	if err := s.repo.Update(ctx, p); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("update patient: %w", err)
	}
	return p, nil
}

func (s *PatientService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.SoftDelete(ctx, id)
}
