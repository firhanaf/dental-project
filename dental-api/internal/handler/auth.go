package handler

import (
	"encoding/json"
	"net/http"

	"github.com/yourusername/dental-api/internal/middleware"
	"github.com/yourusername/dental-api/internal/service"
	"github.com/yourusername/dental-api/pkg/response"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

// POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}
	if req.Email == "" || req.Password == "" {
		response.BadRequest(w, "Email dan password wajib diisi")
		return
	}

	result, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		// Jangan kasih tahu apakah email atau password yang salah
		response.Error(w, 401, "INVALID_CREDENTIALS", "Email atau password salah")
		return
	}

	response.JSON(w, http.StatusOK, result)
}

// GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	if claims == nil {
		response.Unauthorized(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]interface{}{
		"user_id":   claims.UserID,
		"name":      claims.Name,
		"role":      claims.Role,
		"branch_id": claims.BranchID,
	})
}

// POST /api/v1/auth/logout
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// JWT adalah stateless — logout cukup hapus token di client
	// Jika butuh blacklist token, implementasi di sini dengan Redis
	response.JSON(w, http.StatusOK, map[string]string{"message": "Logout berhasil"})
}
