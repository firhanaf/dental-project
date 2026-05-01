package handler

import (
	"net/http"

	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type BranchHandler struct{ svc *service.BranchService }

func NewBranchHandler(svc *service.BranchService) *BranchHandler {
	return &BranchHandler{svc: svc}
}

func (h *BranchHandler) List(w http.ResponseWriter, r *http.Request) {
	branches, err := h.svc.List(r.Context())
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, branches)
}
