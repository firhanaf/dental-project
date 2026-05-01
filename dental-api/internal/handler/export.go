package handler

import (
	"net/http"
	"time"

	"github.com/yourusername/dental-api/internal/dto"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type ExportHandler struct{ svc *service.ExportService }

func NewExportHandler(svc *service.ExportService) *ExportHandler {
	return &ExportHandler{svc: svc}
}

func (h *ExportHandler) ExportPatients(w http.ResponseWriter, r *http.Request) {
	q := dto.ExportQuery{BranchID: r.URL.Query().Get("branch_id")}

	buf, err := h.svc.ExportPatients(r.Context(), q)
	if err != nil {
		response.InternalError(w)
		return
	}

	filename := "export_pasien_" + time.Now().Format("2006-01-02") + ".xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}

func (h *ExportHandler) ExportVisits(w http.ResponseWriter, r *http.Request) {
	q := dto.ExportQuery{
		BranchID: r.URL.Query().Get("branch_id"),
		DateFrom: r.URL.Query().Get("date_from"),
		DateTo:   r.URL.Query().Get("date_to"),
	}

	buf, err := h.svc.ExportVisits(r.Context(), q)
	if err != nil {
		response.InternalError(w)
		return
	}

	filename := "export_kunjungan_" + time.Now().Format("2006-01-02") + ".xlsx"
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(buf.Bytes())
}
