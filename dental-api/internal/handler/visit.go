package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/middleware"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type VisitHandler struct{ svc *service.VisitService }

func NewVisitHandler(svc *service.VisitService) *VisitHandler {
	return &VisitHandler{svc: svc}
}

func (h *VisitHandler) ListByPatient(w http.ResponseWriter, r *http.Request) {
	patientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID pasien tidak valid")
		return
	}

	visits, err := h.svc.ListByPatient(r.Context(), patientID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, visits)
}

func (h *VisitHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	visit, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, visit)
}

func (h *VisitHandler) Create(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	var req dto.CreateVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	visit, err := h.svc.Create(r.Context(), &req, claims)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, visit)
}

func (h *VisitHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	var req dto.UpdateVisitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	visit, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, visit)
}

func (h *VisitHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.JSON(w, http.StatusOK, map[string]string{"message": "Kunjungan berhasil dihapus"})
}
