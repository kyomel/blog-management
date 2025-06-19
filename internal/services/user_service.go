package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/repositories"
)

// UserService handles user-related business logic
type UserService struct {
	repo repositories.UserRepository
}

// NewUserService creates a new instance of UserService
func NewUserService(repo repositories.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// UpdateAvatarURL updates the avatar URL for a user
func (s *UserService) UpdateAvatarURL(ctx context.Context, userID string, avatarURL string) error {
	// Parse the user ID from string to UUID
	id, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	// Call the repository to update the avatar URL
	return s.repo.UpdateAvatarURL(ctx, id, avatarURL)
}
