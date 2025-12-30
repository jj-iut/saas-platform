package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/yourcompany/saas-platform/internal/config"
)

type Service struct {
	repo     *Repository
	jwtConfig config.JWTConfig
}

func NewService(repo *Repository, jwtConfig config.JWTConfig) *Service {
	return &Service{
		repo:      repo,
		jwtConfig: jwtConfig,
	}
}

func (s *Service) Register(req *RegisterRequest) (*AuthResponse, error) {
	// Check if user already exists
	existingUser, _ := s.repo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &User{
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         RoleUser,
	}
	if req.Name != "" {
		user.Name = &req.Name
	}

	if err := s.repo.CreateUser(user); err != nil {
		return nil, err
	}

	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtConfig.AccessSecret, s.jwtConfig.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(user.ID, user.Email, s.jwtConfig.RefreshSecret, s.jwtConfig.RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(req *LoginRequest) (*AuthResponse, error) {
	// Get user by email
	user, err := s.repo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtConfig.AccessSecret, s.jwtConfig.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := GenerateRefreshToken(user.ID, user.Email, s.jwtConfig.RefreshSecret, s.jwtConfig.RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) RefreshToken(refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	claims, err := ValidateToken(refreshToken, s.jwtConfig.RefreshSecret)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.repo.GetUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Generate new tokens
	accessToken, err := GenerateAccessToken(user.ID, user.Email, user.Role, s.jwtConfig.AccessSecret, s.jwtConfig.AccessTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	newRefreshToken, err := GenerateRefreshToken(user.ID, user.Email, s.jwtConfig.RefreshSecret, s.jwtConfig.RefreshTTL)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &AuthResponse{
		User:         user,
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *Service) GetUserByID(userID int64) (*User, error) {
	return s.repo.GetUserByID(userID)
}
