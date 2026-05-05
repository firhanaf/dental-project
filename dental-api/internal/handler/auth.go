package handler

import (
	"encoding/json"
	"errors"
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
		if errors.Is(err, service.ErrAkunTidakAktif) {
			response.Error(w, 403, "ACCOUNT_INACTIVE", "Akun Anda tidak aktif. Hubungi admin untuk mengaktifkan kembali.")
			return
		}
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

// PUT /api/v1/auth/change-password
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	claims := middleware.GetClaims(r)
	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		response.BadRequest(w, "Password saat ini dan password baru wajib diisi")
		return
	}
	err := h.svc.ChangePassword(r.Context(), claims.UserID, req.CurrentPassword, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrPasswordSalah) {
			response.Error(w, http.StatusUnprocessableEntity, "WRONG_PASSWORD", err.Error())
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Password berhasil diubah"})
}

// POST /api/v1/auth/reset-password
func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "Format request tidak valid")
		return
	}
	if req.Email == "" || req.Token == "" || req.NewPassword == "" {
		response.BadRequest(w, "Email, kode reset, dan password baru wajib diisi")
		return
	}
	err := h.svc.ResetPassword(r.Context(), req.Email, req.Token, req.NewPassword)
	if err != nil {
		if errors.Is(err, service.ErrTokenInvalid) {
			response.Error(w, http.StatusUnprocessableEntity, "TOKEN_INVALID", err.Error())
			return
		}
		response.BadRequest(w, err.Error())
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "Password berhasil direset. Silakan login."})
}
