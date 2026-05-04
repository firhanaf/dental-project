package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type UserMgmtHandler struct{ svc *service.UserMgmtService }

func NewUserMgmtHandler(svc *service.UserMgmtService) *UserMgmtHandler {
	return &UserMgmtHandler{svc: svc}
}

func (h *UserMgmtHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.svc.List(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, users)
}

func (h *UserMgmtHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	user, err := h.svc.Create(r.Context(), &req)
	if err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusCreated, user)
}

func (h *UserMgmtHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}

	user, err := h.svc.Update(r.Context(), id, &req)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			response.NotFound(w)
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, user)
}

func (h *UserMgmtHandler) Deactivate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}
	if err := h.svc.Deactivate(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "User berhasil dinonaktifkan"})
}

func (h *UserMgmtHandler) Activate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}
	if err := h.svc.Activate(r.Context(), id); err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "User berhasil diaktifkan"})
}

func (h *UserMgmtHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.BadRequest(w, "ID tidak valid")
		return
	}
	if err := h.svc.Delete(r.Context(), id); err != nil {
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "User berhasil dihapus"})
}
