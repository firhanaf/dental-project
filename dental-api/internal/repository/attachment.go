package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type AttachmentRepo struct{ db *pgxpool.Pool }

func NewAttachmentRepo(db *pgxpool.Pool) *AttachmentRepo { return &AttachmentRepo{db: db} }

func (r *AttachmentRepo) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]model.Attachment, error) {
	rows, err := r.db.Query(ctx, `
		SELECT a.id, a.visit_id, a.uploaded_by, a.original_name, a.stored_name,
		       a.file_path, a.file_type, a.mime_type, a.size_bytes, a.deleted_at, a.created_at
		FROM attachments a
		JOIN visits v ON v.id = a.visit_id
		WHERE v.patient_id = $1 AND a.deleted_at IS NULL
		ORDER BY a.created_at DESC`, patientID)
	if err != nil {
		return nil, fmt.Errorf("list attachments: %w", err)
	}
	defer rows.Close()

	attachments := make([]model.Attachment, 0)
	for rows.Next() {
		var a model.Attachment
		if err := rows.Scan(
			&a.ID, &a.VisitID, &a.UploadedBy, &a.OriginalName, &a.StoredName,
			&a.FilePath, &a.FileType, &a.MimeType, &a.SizeBytes, &a.DeletedAt, &a.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan attachment: %w", err)
		}
		attachments = append(attachments, a)
	}
	return attachments, nil
}

func (r *AttachmentRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	var a model.Attachment
	err := r.db.QueryRow(ctx, `
		SELECT id, visit_id, uploaded_by, original_name, stored_name,
		       file_path, file_type, mime_type, size_bytes, deleted_at, created_at
		FROM attachments
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&a.ID, &a.VisitID, &a.UploadedBy, &a.OriginalName, &a.StoredName,
		&a.FilePath, &a.FileType, &a.MimeType, &a.SizeBytes, &a.DeletedAt, &a.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get attachment: %w", err)
	}
	return &a, nil
}

func (r *AttachmentRepo) Create(ctx context.Context, a *model.Attachment) error {
	a.ID = uuid.New()
	return r.db.QueryRow(ctx, `
		INSERT INTO attachments
		  (id, visit_id, uploaded_by, original_name, stored_name, file_path, file_type, mime_type, size_bytes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)
		RETURNING created_at`,
		a.ID, a.VisitID, a.UploadedBy, a.OriginalName, a.StoredName,
		a.FilePath, a.FileType, a.MimeType, a.SizeBytes,
	).Scan(&a.CreatedAt)
}

// SoftDelete mengembalikan attachment (beserta FilePath) sebelum dihapus agar file fisik bisa ikut dihapus.
func (r *AttachmentRepo) SoftDelete(ctx context.Context, id uuid.UUID) (*model.Attachment, error) {
	a, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	result, err := r.db.Exec(ctx,
		"UPDATE attachments SET deleted_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return nil, fmt.Errorf("soft delete attachment: %w", err)
	}
	if result.RowsAffected() == 0 {
		return nil, fmt.Errorf("record not found")
	}
	return a, nil
}
