package service

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
	"github.com/yourusername/dental-api/pkg/storage"
)

type AttachmentService struct {
	repo      *repository.AttachmentRepo
	visitRepo *repository.VisitRepo
	store     *storage.Local
}

func NewAttachmentService(repo *repository.AttachmentRepo, visitRepo *repository.VisitRepo, store *storage.Local) *AttachmentService {
	return &AttachmentService{repo: repo, visitRepo: visitRepo, store: store}
}

func (s *AttachmentService) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]model.Attachment, error) {
	return s.repo.ListByPatient(ctx, patientID)
}

func (s *AttachmentService) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (s *AttachmentService) Upload(
	ctx context.Context,
	visitIDStr string,
	file multipart.File,
	header *multipart.FileHeader,
	claims *jwtpkg.Claims,
) (*model.Attachment, error) {
	visitID, err := uuid.Parse(visitIDStr)
	if err != nil {
		return nil, fmt.Errorf("visit_id tidak valid")
	}

	// Ambil patientID dari visit untuk path penyimpanan file
	visit, err := s.visitRepo.GetByID(ctx, visitID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("kunjungan tidak ditemukan")
		}
		return nil, fmt.Errorf("cari kunjungan: %w", err)
	}

	storedName, filePath, mimeType, err := s.store.Save(file, header, visit.PatientID.String(), visitIDStr)
	if err != nil {
		return nil, err
	}

	fileType := model.FileTypeImage
	if mimeType == "application/pdf" {
		fileType = model.FileTypePDF
	}

	a := &model.Attachment{
		VisitID:      visitID,
		UploadedBy:   claims.UserID,
		OriginalName: header.Filename,
		StoredName:   storedName,
		FilePath:     filePath,
		FileType:     fileType,
		MimeType:     mimeType,
		SizeBytes:    header.Size,
	}

	if err := s.repo.Create(ctx, a); err != nil {
		_ = s.store.Delete(filePath)
		return nil, fmt.Errorf("simpan attachment: %w", err)
	}
	return a, nil
}

func (s *AttachmentService) Download(ctx context.Context, id uuid.UUID, w http.ResponseWriter) error {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrNotFound
		}
		return err
	}
	return s.store.ServeFile(w, a.FilePath, a.MimeType, a.OriginalName)
}

func (s *AttachmentService) Delete(ctx context.Context, id uuid.UUID) error {
	a, err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || err.Error() == "record not found" {
			return ErrNotFound
		}
		return err
	}
	return s.store.Delete(a.FilePath)
}
