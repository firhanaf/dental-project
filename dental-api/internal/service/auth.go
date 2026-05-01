package service

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/yourusername/dental-api/internal/repository"
	jwtpkg "github.com/yourusername/dental-api/pkg/jwt"
)

type AuthService struct {
	userRepo *repository.UserRepo
	jwt      *jwtpkg.Manager
}

func NewAuthService(userRepo *repository.UserRepo, jwt *jwtpkg.Manager) *AuthService {
	return &AuthService{userRepo: userRepo, jwt: jwt}
}

type LoginResult struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	User      struct {
		ID       string `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		BranchID *string `json:"branch_id"`
	} `json:"user"`
}

func (s *AuthService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}
	if !user.IsActive {
		return nil, errors.New("akun tidak aktif")
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
	result.User.ID    = user.ID.String()
	result.User.Name  = user.Name
	result.User.Email = user.Email
	result.User.Role  = string(user.Role)
	if user.BranchID != nil {
		s := user.BranchID.String()
		result.User.BranchID = &s
	}

	return result, nil
}
