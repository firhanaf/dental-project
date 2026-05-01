package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/yourusername/dental-api/internal/model"
)

type PatientRepo struct {
	db *pgxpool.Pool
}

func NewPatientRepo(db *pgxpool.Pool) *PatientRepo {
	return &PatientRepo{db: db}
}

type PatientFilter struct {
	BranchID *uuid.UUID
	Status   *model.PatientStatus
	Search   string // nama, no_rm, atau phone
	Page     int
	Limit    int
}

func (r *PatientRepo) List(ctx context.Context, f PatientFilter) ([]model.PatientListRow, int, error) {
	if f.Limit == 0 {
		f.Limit = 20
	}
	if f.Page < 1 {
		f.Page = 1
	}
	offset := (f.Page - 1) * f.Limit

	// Query dinamis dengan kondisi opsional
	where := "WHERE 1=1"
	args  := []interface{}{}
	i     := 1

	if f.BranchID != nil {
		where += fmt.Sprintf(" AND branch_id = $%d", i)
		args = append(args, *f.BranchID)
		i++
	}
	if f.Status != nil {
		where += fmt.Sprintf(" AND status = $%d", i)
		args = append(args, *f.Status)
		i++
	}
	if f.Search != "" {
		// Cari di nama (trigram), no_rm, dan phone
		where += fmt.Sprintf(
			" AND (name ILIKE $%d OR no_rm ILIKE $%d OR phone ILIKE $%d)",
			i, i, i,
		)
		args = append(args, "%"+f.Search+"%")
		i++
	}

	// Hitung total (untuk pagination)
	var total int
	countQ := "SELECT COUNT(*) FROM v_patient_list " + where
	if err := r.db.QueryRow(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count patients: %w", err)
	}

	// Query data dengan LIMIT + OFFSET
	query := fmt.Sprintf(`
		SELECT id, branch_id, branch_name, no_rm, name, date_of_birth, age,
		       gender, phone, allergy_notes, status, created_at, updated_at,
		       last_visit_date, last_diagnosis, last_doctor, total_visits, total_cost
		FROM v_patient_list %s
		ORDER BY updated_at DESC
		LIMIT $%d OFFSET $%d
	`, where, i, i+1)
	args = append(args, f.Limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("list patients: %w", err)
	}
	defer rows.Close()

	patients := make([]model.PatientListRow, 0)
	for rows.Next() {
		var p model.PatientListRow
		if err := rows.Scan(
			&p.ID, &p.BranchID, &p.BranchName, &p.NoRM, &p.Name,
			&p.DateOfBirth, &p.Age, &p.Gender, &p.Phone, &p.AllergyNotes,
			&p.Status, &p.CreatedAt, &p.UpdatedAt,
			&p.LastVisitDate, &p.LastDiagnosis, &p.LastDoctor,
			&p.TotalVisits, &p.TotalCost,
		); err != nil {
			return nil, 0, fmt.Errorf("scan patient row: %w", err)
		}
		patients = append(patients, p)
	}

	return patients, total, nil
}

func (r *PatientRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Patient, error) {
	query := `
		SELECT id, branch_id, created_by, no_rm, name, nik, date_of_birth,
		       gender, phone, address, occupation, allergy_notes, status,
		       deleted_at, created_at, updated_at
		FROM patients
		WHERE id = $1 AND deleted_at IS NULL
	`
	var p model.Patient
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID, &p.BranchID, &p.CreatedBy, &p.NoRM, &p.Name, &p.NIK,
		&p.DateOfBirth, &p.Gender, &p.Phone, &p.Address, &p.Occupation,
		&p.AllergyNotes, &p.Status, &p.DeletedAt, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("get patient: %w", err)
	}
	return &p, nil
}

func (r *PatientRepo) Create(ctx context.Context, p *model.Patient) error {
	// Generate No. RM via PostgreSQL function
	var noRM string
	if err := r.db.QueryRow(ctx, "SELECT generate_no_rm()").Scan(&noRM); err != nil {
		return fmt.Errorf("generate no_rm: %w", err)
	}
	p.NoRM = noRM
	p.ID   = uuid.New()

	query := `
		INSERT INTO patients
		  (id, branch_id, created_by, no_rm, name, nik, date_of_birth,
		   gender, phone, address, occupation, allergy_notes)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
		RETURNING status, created_at, updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.ID, p.BranchID, p.CreatedBy, p.NoRM, p.Name, p.NIK,
		p.DateOfBirth, p.Gender, p.Phone, p.Address, p.Occupation, p.AllergyNotes,
	).Scan(&p.Status, &p.CreatedAt, &p.UpdatedAt)
}

func (r *PatientRepo) Update(ctx context.Context, p *model.Patient) error {
	query := `
		UPDATE patients SET
			name=$1, nik=$2, date_of_birth=$3, gender=$4, phone=$5,
			address=$6, occupation=$7, allergy_notes=$8, updated_at=NOW()
		WHERE id=$9 AND deleted_at IS NULL
		RETURNING updated_at
	`
	return r.db.QueryRow(ctx, query,
		p.Name, p.NIK, p.DateOfBirth, p.Gender, p.Phone,
		p.Address, p.Occupation, p.AllergyNotes, p.ID,
	).Scan(&p.UpdatedAt)
}

func (r *PatientRepo) ExistsByNIK(ctx context.Context, nik string, excludeID *uuid.UUID) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM patients WHERE nik = $1 AND deleted_at IS NULL"
	args := []interface{}{nik}
	if excludeID != nil {
		query += " AND id != $2"
		args = append(args, *excludeID)
	}
	query += ")"
	var exists bool
	if err := r.db.QueryRow(ctx, query, args...).Scan(&exists); err != nil {
		return false, fmt.Errorf("check nik: %w", err)
	}
	return exists, nil
}

func (r *PatientRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		"UPDATE patients SET deleted_at=NOW(), updated_at=NOW() WHERE id=$1 AND deleted_at IS NULL",
		id,
	)
	return err
}

func (r *PatientRepo) ListForExport(ctx context.Context, branchID *uuid.UUID) ([]model.PatientExportRow, error) {
	where := "WHERE p.deleted_at IS NULL"
	args := []interface{}{}
	if branchID != nil {
		where += " AND p.branch_id = $1"
		args = append(args, *branchID)
	}

	query := fmt.Sprintf(`
		SELECT p.no_rm, p.name, p.date_of_birth,
		       EXTRACT(YEAR FROM AGE(p.date_of_birth))::INT AS age,
		       p.gender, p.phone, p.address, p.allergy_notes, p.status,
		       b.name AS branch_name, p.created_at
		FROM patients p
		JOIN branches b ON b.id = p.branch_id
		%s
		ORDER BY p.created_at DESC`, where)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("list patients for export: %w", err)
	}
	defer rows.Close()

	result := make([]model.PatientExportRow, 0)
	for rows.Next() {
		var p model.PatientExportRow
		if err := rows.Scan(
			&p.NoRM, &p.Name, &p.DateOfBirth, &p.Age,
			&p.Gender, &p.Phone, &p.Address, &p.AllergyNotes,
			&p.Status, &p.BranchName, &p.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan patient export: %w", err)
		}
		result = append(result, p)
	}
	return result, nil
}
