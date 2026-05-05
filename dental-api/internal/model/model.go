package model

import (
	"time"
	"github.com/google/uuid"
)

// ── Enums ────────────────────────────────────────────────────

type UserRole      string
type GenderType    string
type PatientStatus string
type FileType      string
type AuditAction   string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleWrite      UserRole = "write"
	RoleReadonly   UserRole = "readonly"

	GenderMale   GenderType = "male"
	GenderFemale GenderType = "female"

	StatusNew          PatientStatus = "new"
	StatusActive       PatientStatus = "active"
	StatusNeedsControl PatientStatus = "needs_control"

	FileTypePDF   FileType = "pdf"
	FileTypeImage FileType = "image"

	AuditCreate  AuditAction = "create"
	AuditUpdate  AuditAction = "update"
	AuditDelete  AuditAction = "delete"
	AuditRestore AuditAction = "restore"
)

// ── Domain Models ─────────────────────────────────────────────

type Branch struct {
	ID        uuid.UUID `db:"id"         json:"id"`
	Name      string    `db:"name"       json:"name"`
	Address   *string   `db:"address"    json:"address"`
	Phone     *string   `db:"phone"      json:"phone"`
	IsActive  bool      `db:"is_active"  json:"is_active"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}

type User struct {
	ID           uuid.UUID  `db:"id"             json:"id"`
	BranchID     *uuid.UUID `db:"branch_id"      json:"branch_id"`
	Name         string     `db:"name"           json:"name"`
	Email        string     `db:"email"          json:"email"`
	PasswordHash string     `db:"password_hash"  json:"-"` // tidak pernah dikirim ke client
	Role         UserRole   `db:"role"           json:"role"`
	IsActive     bool       `db:"is_active"      json:"is_active"`
	LastLoginAt  *time.Time `db:"last_login_at"  json:"last_login_at"`
	CreatedAt    time.Time  `db:"created_at"     json:"created_at"`
	UpdatedAt    time.Time  `db:"updated_at"     json:"updated_at"`
}

type Patient struct {
	ID           uuid.UUID     `db:"id"            json:"id"`
	BranchID     uuid.UUID     `db:"branch_id"     json:"branch_id"`
	CreatedBy    uuid.UUID     `db:"created_by"    json:"created_by"`
	NoRM         string        `db:"no_rm"         json:"no_rm"`
	Name         string        `db:"name"          json:"name"`
	NIK          *string       `db:"nik"           json:"nik"`
	DateOfBirth  time.Time     `db:"date_of_birth" json:"date_of_birth"`
	Gender       GenderType    `db:"gender"        json:"gender"`
	Phone        string        `db:"phone"         json:"phone"`
	Address      *string       `db:"address"       json:"address"`
	Occupation   *string       `db:"occupation"    json:"occupation"`
	AllergyNotes *string       `db:"allergy_notes" json:"allergy_notes"`
	Status       PatientStatus `db:"status"        json:"status"`
	DeletedAt    *time.Time    `db:"deleted_at"    json:"deleted_at,omitempty"`
	CreatedAt    time.Time     `db:"created_at"    json:"created_at"`
	UpdatedAt    time.Time     `db:"updated_at"    json:"updated_at"`
}

type Visit struct {
	ID              uuid.UUID  `db:"id"               json:"id"`
	PatientID       uuid.UUID  `db:"patient_id"       json:"patient_id"`
	BranchID        uuid.UUID  `db:"branch_id"        json:"branch_id"`
	DoctorID        uuid.UUID  `db:"doctor_id"        json:"doctor_id"`
	CreatedBy       uuid.UUID  `db:"created_by"       json:"created_by"`
	VisitDate       time.Time  `db:"visit_date"       json:"visit_date"`
	ChiefComplaint  string     `db:"chief_complaint"  json:"chief_complaint"`
	Diagnosis       *string    `db:"diagnosis"        json:"diagnosis"`
	Treatment       *string    `db:"treatment"        json:"treatment"`
	TeethInvolved   *string    `db:"teeth_involved"   json:"teeth_involved"`
	Cost            float64    `db:"cost"             json:"cost"`
	NextControlDate *time.Time `db:"next_control_date" json:"next_control_date"`
	Notes           *string    `db:"notes"            json:"notes"`
	DeletedAt       *time.Time `db:"deleted_at"       json:"deleted_at,omitempty"`
	CreatedAt       time.Time  `db:"created_at"       json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at"       json:"updated_at"`
}

type Attachment struct {
	ID           uuid.UUID  `db:"id"            json:"id"`
	VisitID      uuid.UUID  `db:"visit_id"      json:"visit_id"`
	UploadedBy   uuid.UUID  `db:"uploaded_by"   json:"uploaded_by"`
	OriginalName string     `db:"original_name" json:"original_name"`
	StoredName   string     `db:"stored_name"   json:"-"`
	FilePath     string     `db:"file_path"     json:"-"`
	FileType     FileType   `db:"file_type"     json:"file_type"`
	MimeType     string     `db:"mime_type"     json:"mime_type"`
	SizeBytes    int64      `db:"size_bytes"    json:"size_bytes"`
	DeletedAt    *time.Time `db:"deleted_at"    json:"deleted_at,omitempty"`
	CreatedAt    time.Time  `db:"created_at"    json:"created_at"`
}

type AuditLog struct {
	ID        int64       `db:"id"         json:"id"`
	UserID    uuid.UUID   `db:"user_id"    json:"user_id"`
	Action    AuditAction `db:"action"     json:"action"`
	TableName string      `db:"table_name" json:"table_name"`
	RecordID  uuid.UUID   `db:"record_id"  json:"record_id"`
	OldData   *string     `db:"old_data"   json:"old_data"`
	NewData   *string     `db:"new_data"   json:"new_data"`
	IPAddress *string     `db:"ip_address" json:"ip_address"`
	CreatedAt time.Time   `db:"created_at" json:"created_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID  `db:"id"`
	UserID    uuid.UUID  `db:"user_id"`
	TokenHash string     `db:"token_hash"`
	ExpiresAt time.Time  `db:"expires_at"`
	UsedAt    *time.Time `db:"used_at"`
	CreatedAt time.Time  `db:"created_at"`
}

// ── View Models (hasil JOIN, hanya untuk read) ────────────────

type PatientListRow struct {
	Patient
	BranchName    string  `db:"branch_name"    json:"branch_name"`
	Age           int     `db:"age"            json:"age"`
	LastVisitDate *time.Time `db:"last_visit_date" json:"last_visit_date"`
	LastDiagnosis *string `db:"last_diagnosis" json:"last_diagnosis"`
	LastDoctor    *string `db:"last_doctor"    json:"last_doctor"`
	TotalVisits   int     `db:"total_visits"   json:"total_visits"`
	TotalCost     float64 `db:"total_cost"     json:"total_cost"`
}

// PatientExportRow — data pasien lengkap untuk export Excel
type PatientExportRow struct {
	NoRM         string        `db:"no_rm"`
	Name         string        `db:"name"`
	DateOfBirth  time.Time     `db:"date_of_birth"`
	Age          int           `db:"age"`
	Gender       GenderType    `db:"gender"`
	Phone        string        `db:"phone"`
	Address      *string       `db:"address"`
	AllergyNotes *string       `db:"allergy_notes"`
	Status       PatientStatus `db:"status"`
	BranchName   string        `db:"branch_name"`
	CreatedAt    time.Time     `db:"created_at"`
}

// VisitExportRow — data kunjungan + join untuk export Excel
type VisitExportRow struct {
	PatientNoRM    string    `db:"no_rm"`
	PatientName    string    `db:"patient_name"`
	VisitDate      time.Time `db:"visit_date"`
	DoctorName     string    `db:"doctor_name"`
	ChiefComplaint string    `db:"chief_complaint"`
	Diagnosis      *string   `db:"diagnosis"`
	Treatment      *string   `db:"treatment"`
	TeethInvolved  *string   `db:"teeth_involved"`
	Cost           float64   `db:"cost"`
	BranchName     string    `db:"branch_name"`
}
