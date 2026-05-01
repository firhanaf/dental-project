package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type VisitRepo struct{ db *pgxpool.Pool }

func NewVisitRepo(db *pgxpool.Pool) *VisitRepo { return &VisitRepo{db: db} }

type ExportVisitFilter struct {
	BranchID *uuid.UUID
	DateFrom *time.Time
	DateTo   *time.Time
}

func (r *VisitRepo) ListByPatient(ctx context.Context, patientID uuid.UUID) ([]model.Visit, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id, patient_id, branch_id, doctor_id, created_by, visit_date, chief_complaint,
		       diagnosis, treatment, teeth_involved, cost, next_control_date, notes,
		       deleted_at, created_at, updated_at
		FROM visits
		WHERE patient_id = $1 AND deleted_at IS NULL
		ORDER BY visit_date DESC`, patientID)
	if err != nil {
		return nil, fmt.Errorf("list visits: %w", err)
	}
	defer rows.Close()

	visits := make([]model.Visit, 0)
	for rows.Next() {
		var v model.Visit
		if err := rows.Scan(
			&v.ID, &v.PatientID, &v.BranchID, &v.DoctorID, &v.CreatedBy,
			&v.VisitDate, &v.ChiefComplaint, &v.Diagnosis, &v.Treatment,
			&v.TeethInvolved, &v.Cost, &v.NextControlDate, &v.Notes,
			&v.DeletedAt, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan visit: %w", err)
		}
		visits = append(visits, v)
	}
	return visits, nil
}

func (r *VisitRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Visit, error) {
	var v model.Visit
	err := r.db.QueryRow(ctx, `
		SELECT id, patient_id, branch_id, doctor_id, created_by, visit_date, chief_complaint,
		       diagnosis, treatment, teeth_involved, cost, next_control_date, notes,
		       deleted_at, created_at, updated_at
		FROM visits
		WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&v.ID, &v.PatientID, &v.BranchID, &v.DoctorID, &v.CreatedBy,
		&v.VisitDate, &v.ChiefComplaint, &v.Diagnosis, &v.Treatment,
		&v.TeethInvolved, &v.Cost, &v.NextControlDate, &v.Notes,
		&v.DeletedAt, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get visit: %w", err)
	}
	return &v, nil
}

func (r *VisitRepo) Create(ctx context.Context, v *model.Visit) error {
	v.ID = uuid.New()
	return r.db.QueryRow(ctx, `
		INSERT INTO visits
		  (id, patient_id, branch_id, doctor_id, created_by, visit_date, chief_complaint,
		   diagnosis, treatment, teeth_involved, cost, next_control_date, notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)
		RETURNING created_at, updated_at`,
		v.ID, v.PatientID, v.BranchID, v.DoctorID, v.CreatedBy, v.VisitDate,
		v.ChiefComplaint, v.Diagnosis, v.Treatment, v.TeethInvolved, v.Cost,
		v.NextControlDate, v.Notes,
	).Scan(&v.CreatedAt, &v.UpdatedAt)
}

func (r *VisitRepo) Update(ctx context.Context, v *model.Visit) error {
	return r.db.QueryRow(ctx, `
		UPDATE visits SET
		    doctor_id=$1, visit_date=$2, chief_complaint=$3, diagnosis=$4,
		    treatment=$5, teeth_involved=$6, cost=$7, next_control_date=$8,
		    notes=$9, updated_at=NOW()
		WHERE id=$10 AND deleted_at IS NULL
		RETURNING updated_at`,
		v.DoctorID, v.VisitDate, v.ChiefComplaint, v.Diagnosis, v.Treatment,
		v.TeethInvolved, v.Cost, v.NextControlDate, v.Notes, v.ID,
	).Scan(&v.UpdatedAt)
}

func (r *VisitRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	result, err := r.db.Exec(ctx,
		"UPDATE visits SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL", id)
	if err != nil {
		return fmt.Errorf("soft delete visit: %w", err)
	}
	if result.RowsAffected() == 0 {
		return fmt.Errorf("record not found")
	}
	return nil
}

func (r *VisitRepo) ListForExport(ctx context.Context, f ExportVisitFilter) ([]model.VisitExportRow, error) {
	where := "WHERE v.deleted_at IS NULL"
	args := []interface{}{}
	i := 1

	if f.BranchID != nil {
		where += fmt.Sprintf(" AND v.branch_id = $%d", i)
		args = append(args, *f.BranchID)
		i++
	}
	if f.DateFrom != nil {
		where += fmt.Sprintf(" AND v.visit_date >= $%d", i)
		args = append(args, *f.DateFrom)
		i++
	}
	if f.DateTo != nil {
		where += fmt.Sprintf(" AND v.visit_date <= $%d", i)
		args = append(args, *f.DateTo)
		i++
	}
	_ = i

	query := fmt.Sprintf(`
		SELECT p.no_rm, p.name AS patient_name, v.visit_date, u.name AS doctor_name,
		       v.chief_complaint, v.diagnosis, v.treatment, v.teeth_involved, v.cost, b.name AS branch_name
		FROM visits v
		JOIN patients p ON p.id = v.patient_id
		JOIN users u ON u.id = v.doctor_id
		JOIN branches b ON b.id = v.branch_id
		%s ORDER BY v.visit_date DESC`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list visits for export: %w", err)
	}
	defer rows.Close()

	result := make([]model.VisitExportRow, 0)
	for rows.Next() {
		var row model.VisitExportRow
		if err := rows.Scan(
			&row.PatientNoRM, &row.PatientName, &row.VisitDate, &row.DoctorName,
			&row.ChiefComplaint, &row.Diagnosis, &row.Treatment, &row.TeethInvolved,
			&row.Cost, &row.BranchName,
		); err != nil {
			return nil, fmt.Errorf("scan visit export: %w", err)
		}
		result = append(result, row)
	}
	return result, nil
}
