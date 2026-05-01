package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/middleware"
	"github.com/yourusername/dental-api/internal/model"
	"github.com/yourusername/dental-api/internal/repository"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type PatientHandler struct{ svc *service.PatientService }

func NewPatientHandler(svc *service.PatientService) *PatientHandler {
	return &PatientHandler{svc: svc}
}

func (h *PatientHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	filter := repository.PatientFilter{
		Search: r.URL.Query().Get("search"),
		Page:   page,
		Limit:  limit,
	}

	if statusStr := r.URL.Query().Get("status"); statusStr != "" {
		s := model.PatientStatus(statusStr)
		filter.Status = &s
	}

	if branchIDStr := r.URL.Query().Get("branch_id"); branchIDStr != "" {
		id, err := uuid.Parse(branchIDStr)
		if err != nil {
			response.BadRequest(w, "branch_id tidak valid")
			return
		}
		filter.BranchID = &id
	}

	patients, total, err := h.svc.List(r.Context(), filter)
	if err != nil {
		response.InternalError(w)
		return
	}

	hasNext := (page * limit) < total
	response.WithMeta(w, http.StatusOK, patients, &response.Meta{
		Page: page, Limit: limit, Total: total, HasNext: hasNext,
	})
}

func (h *PatientHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	patient, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, patient)
}

func (h *PatientHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	var req dto.CreatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	patient, err := h.svc.Create(r.Context(), &req, claims)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, patient)
}

func (h *PatientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	var req dto.UpdatePatientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	patient, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, patient)
}

func (h *PatientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	if err := h.svc.Delete(r.Context(), id); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Pasien berhasil dihapus"})
}
