package response

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorBody  `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Meta struct {
	Page    int `json:"page"`
	Limit   int `json:"limit"`
	Total   int `json:"total"`
	HasNext bool `json:"has_next"`
}

func JSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Success: true, Data: data})
}

func WithMeta(w http.ResponseWriter, status int, data interface{}, meta *Meta) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{Success: true, Data: data, Meta: meta})
}

func Error(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Envelope{
		Success: false,
		Error:   &ErrorBody{Code: code, Message: message},
	})
}

func Unauthorized(w http.ResponseWriter)  { Error(w, 401, "UNAUTHORIZED", "Login diperlukan") }
func Forbidden(w http.ResponseWriter)     { Error(w, 403, "FORBIDDEN", "Akses ditolak") }
func NotFound(w http.ResponseWriter)      { Error(w, 404, "NOT_FOUND", "Data tidak ditemukan") }
func BadRequest(w http.ResponseWriter, msg string) { Error(w, 400, "BAD_REQUEST", msg) }
func InternalError(w http.ResponseWriter) { Error(w, 500, "INTERNAL_ERROR", "Terjadi kesalahan server") }
