package dto

import "time"

type CreatePatientRequest struct {
	BranchID     string  `json:"branch_id"`
	Name         string  `json:"name"`
	NIK          *string `json:"nik"`
	DateOfBirth  string  `json:"date_of_birth"` // "YYYY-MM-DD"
	Gender       string  `json:"gender"`
	Phone        string  `json:"phone"`
	Address      *string `json:"address"`
	Occupation   *string `json:"occupation"`
	AllergyNotes *string `json:"allergy_notes"`
}

type UpdatePatientRequest struct {
	Name         string  `json:"name"`
	NIK          *string `json:"nik"`
	DateOfBirth  string  `json:"date_of_birth"`
	Gender       string  `json:"gender"`
	Phone        string  `json:"phone"`
	Address      *string `json:"address"`
	Occupation   *string `json:"occupation"`
	AllergyNotes *string `json:"allergy_notes"`
}

type PatientListQuery struct {
	Search   string `query:"search"`
	BranchID string `query:"branch_id"`
	Status   string `query:"status"`
	Page     int    `query:"page"`
	Limit    int    `query:"limit"`
}

type CreateVisitRequest struct {
	PatientID       string   `json:"patient_id"`
	BranchID        string   `json:"branch_id"`
	DoctorID        string   `json:"doctor_id"`
	VisitDate       string   `json:"visit_date"`
	ChiefComplaint  string   `json:"chief_complaint"`
	Diagnosis       *string  `json:"diagnosis"`
	Treatment       *string  `json:"treatment"`
	TeethInvolved   *string  `json:"teeth_involved"`
	Cost            float64  `json:"cost"`
	NextControlDate *string  `json:"next_control_date"`
	Notes           *string  `json:"notes"`
}

type UpdateVisitRequest struct {
	DoctorID        string   `json:"doctor_id"`
	VisitDate       string   `json:"visit_date"`
	ChiefComplaint  string   `json:"chief_complaint"`
	Diagnosis       *string  `json:"diagnosis"`
	Treatment       *string  `json:"treatment"`
	TeethInvolved   *string  `json:"teeth_involved"`
	Cost            float64  `json:"cost"`
	NextControlDate *string  `json:"next_control_date"`
	Notes           *string  `json:"notes"`
}

type ExportQuery struct {
	BranchID string `query:"branch_id"`
	DateFrom string `query:"date_from"`
	DateTo   string `query:"date_to"`
}

type CreateUserRequest struct {
	BranchID *string `json:"branch_id"`
	Name     string  `json:"name"`
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Role     string  `json:"role"`
}

type UpdateUserRequest struct {
	BranchID *string `json:"branch_id"`
	Name     string  `json:"name"`
	Role     string  `json:"role"`
}

// Parsed time helper — digunakan di service layer
type ParsedDate = time.Time
