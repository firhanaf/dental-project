package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/yourusername/dental-api/internal/repository"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var ErrAkunTidakAktif = errors.New("akun tidak aktif")
var ErrPasswordSalah   = errors.New("password saat ini tidak sesuai")
var ErrTokenInvalid    = errors.New("kode reset tidak valid atau sudah kedaluwarsa")

type AuthService struct {
	userRepo  *repository.UserRepo
	resetRepo *repository.PasswordResetRepo
	jwt       *jwtpkg.Manager
}

func NewAuthService(userRepo *repository.UserRepo, resetRepo *repository.PasswordResetRepo, jwt *jwtpkg.Manager) *AuthService {
	return &AuthService{userRepo: userRepo, resetRepo: resetRepo, jwt: jwt}
}

type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      struct {
		ID       string  `json:"id"`
		Name     string  `json:"name"`
		Email    string  `json:"email"`
		Role     string  `json:"role"`
		BranchID *string `json:"branch_id"`
	} `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, ErrAkunTidakAktif
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token, err := s.jwt.Generate(user)
	if err != nil {
		return nil, err
	}

	// Update last_login_at (fire and forget)
	go s.userRepo.UpdateLastLogin(context.Background(), user.ID)

	result := &LoginResult{
		Token:     token,
		ExpiresAt: time.Now().Add(8 * time.Hour),
	}
	result.User.ID = user.ID.String()
	result.User.Name = user.Name
	result.User.Email = user.Email
	result.User.Role = string(user.Role)
	if user.BranchID != nil {
		s := user.BranchID.String()
		result.User.BranchID = &s
	}

	return result, nil
}

func (s *AuthService) ChangePassword(ctx context.Context, userID uuid.UUID, currentPassword, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("password baru minimal 8 karakter")
	}
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("find user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return ErrPasswordSalah
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}

func (s *AuthService) ResetPassword(ctx context.Context, email, token, newPassword string) error {
	if len(newPassword) < 8 {
		return fmt.Errorf("password baru minimal 8 karakter")
	}
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Jangan bocorkan apakah email ada atau tidak
		return ErrTokenInvalid
	}
	h := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(h[:])

	record, err := s.resetRepo.FindValid(ctx, user.ID, tokenHash)
	if err != nil {
		return ErrTokenInvalid
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	if err := s.userRepo.UpdatePassword(ctx, user.ID, string(hash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}
	return s.resetRepo.MarkUsed(ctx, record.ID)
}
