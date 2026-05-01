package handler

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/middleware"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type AttachmentHandler struct{ svc *service.AttachmentService }

func NewAttachmentHandler(svc *service.AttachmentService) *AttachmentHandler {
	return &AttachmentHandler{svc: svc}
}

func (h *AttachmentHandler) ListByPatient(w http.ResponseWriter, r *http.Request) {
	patientID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID pasien tidak valid")
		return
	}

	attachments, err := h.svc.ListByPatient(r.Context(), patientID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, attachments)
}

func (h *AttachmentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	attachment, err := h.svc.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, attachment)
}

func (h *AttachmentHandler) Upload(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)

	// Batas parsing 32MB untuk form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		response.BadRequest(w, "Form tidak valid")
		return
	}

	visitIDStr := r.FormValue("visit_id")
	if visitIDStr == "" {
		response.BadRequest(w, "visit_id wajib diisi")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.BadRequest(w, "File wajib disertakan")
		return
	}
	defer file.Close()

	attachment, err := h.svc.Upload(r.Context(), visitIDStr, file, header, claims)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, attachment)
}

func (h *AttachmentHandler) Download(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	if err := h.svc.Download(r.Context(), id, w); err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.InternalError(w)
		return
	}
}

func (h *AttachmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
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
	response.JSON(w, http.StatusOK, map[string]string{"message": "Lampiran berhasil dihapus"})
}
