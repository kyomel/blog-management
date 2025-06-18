package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/repositories"
	"github.com/kyomel/blog-management/internal/utils"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotActive      = errors.New("user account is not active")
)

type AuthService interface {
	Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error)
	ValidateToken(tokenString string) (*utils.JWTClaims, error)
}

type authService struct {
	userRepo     repositories.UserRepository
	jwtService   utils.JWTService
	accessExpiry time.Duration
}

func NewAuthService(
	userRepo repositories.UserRepository,
	jwtService utils.JWTService,
	accessExpiry time.Duration,
) AuthService {
	return &authService{
		userRepo:     userRepo,
		jwtService:   jwtService,
		accessExpiry: accessExpiry,
	}
}

func (s *authService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create new user
	user := &models.User{
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: hashedPassword,
		Role:         models.RoleUser, // Default role
		IsActive:     true,
	}

	// Save user to database
	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	// Generate JWT tokens
	tokens, err := s.jwtService.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
	}, nil
}

func (s *authService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	// Find user by email
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if user is active
	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	// Verify password
	err = utils.VerifyPassword(user.PasswordHash, req.Password)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate JWT tokens
	tokens, err := s.jwtService.GenerateTokenPair(
		user.ID,
		user.Username,
		user.Email,
		string(user.Role),
	)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	claims, err := s.jwtService.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, ErrUserNotActive
	}

	tokens, err := s.jwtService.RefreshTokens(refreshToken)
	if err != nil {
		return nil, err
	}

	return &models.AuthResponse{
		User: models.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			Role:      user.Role,
			AvatarURL: user.AvatarURL,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		},
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    int64(s.accessExpiry.Seconds()),
	}, nil
}

func (s *authService) ValidateToken(tokenString string) (*utils.JWTClaims, error) {
	return s.jwtService.ValidateToken(tokenString)
}
