package middleware

import (
	"context"
	"net/http"
	"strings"

	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"github.com/yourusername/dental-api/pkg/response"
)

type contextKey string
const ClaimsKey contextKey = "claims"

type AuthMiddleware struct {
	jwt *jwtpkg.Manager
}

func NewAuth(jwt *jwtpkg.Manager) *AuthMiddleware {
	return &AuthMiddleware{jwt: jwt}
}

// Authenticate — wajib login, semua role boleh
func (a *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.extractClaims(r)
		if err != nil {
			response.Unauthorized(w)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireWrite — hanya role write & superadmin
func (a *AuthMiddleware) RequireWrite(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.extractClaims(r)
		if err != nil {
			response.Unauthorized(w)
			return
		}
		if claims.Role == "readonly" {
			response.Forbidden(w)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireSuperAdmin — hanya superadmin
func (a *AuthMiddleware) RequireSuperAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims, err := a.extractClaims(r)
		if err != nil {
			response.Unauthorized(w)
			return
		}
		if claims.Role != "superadmin" {
			response.Forbidden(w)
			return
		}
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *AuthMiddleware) extractClaims(r *http.Request) (*jwtpkg.Claims, error) {
	authHeader := r.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return nil, http.ErrNoCookie
	}
	return a.jwt.Verify(strings.TrimPrefix(authHeader, "Bearer "))
}

// GetClaims — helper untuk handler
func GetClaims(r *http.Request) *jwtpkg.Claims {
	claims, _ := r.Context().Value(ClaimsKey).(*jwtpkg.Claims)
	return claims
}
