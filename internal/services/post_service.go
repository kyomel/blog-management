package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/repositories"
)

var (
	ErrPostNotFound     = errors.New("post not found")
	ErrPostSlugConflict = errors.New("post slug already exists")
)

// PostService defines the interface for post-related operations
type PostService interface {
	Create(ctx context.Context, req *models.CreatePostRequest) (*models.PostResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.PostResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.PostResponse, error)
	GetAll(ctx context.Context, filter *models.PostFilter, page, pageSize int) (*models.PaginatedPostResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdatePostRequest) (*models.PostResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Publish(ctx context.Context, id uuid.UUID) (*models.PostResponse, error)
	IncrementViewCount(ctx context.Context, id uuid.UUID) error
}

type postService struct {
	repo *repositories.PostRepository
}

// NewPostService creates a new instance of PostService
func NewPostService(repo *repositories.PostRepository) PostService {
	return &postService{
		repo: repo,
	}
}

// Create creates a new post
func (s *postService) Create(ctx context.Context, req *models.CreatePostRequest) (*models.PostResponse, error) {
	// Check if slug already exists
	existingPost, err := s.repo.GetBySlug(req.Slug)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}
	if existingPost != nil {
		return nil, ErrPostSlugConflict
	}

	// Create new post
	post := &models.Post{
		ID:               uuid.New(),
		AuthorID:         req.AuthorID,
		CategoryID:       req.CategoryID,
		Title:            req.Title,
		Slug:             req.Slug,
		Content:          req.Content,
		Excerpt:          req.Excerpt,
		FeaturedImageURL: req.FeaturedImageURL,
		Status:           req.Status,
		IsFeatured:       req.IsFeatured,
		Metadata:         req.Metadata,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	// Set PublishedAt if status is published
	if req.Status == models.StatusPublished {
		now := time.Now()
		post.PublishedAt = &now
	}

	// Create post and associate tags
	if err := s.repo.Create(post, req.TagIDs); err != nil {
		return nil, err
	}

	// Get the created post with all relationships
	createdPost, err := s.repo.GetByID(post.ID)
	if err != nil {
		return nil, err
	}

	return s.mapPostToResponse(createdPost), nil
}

// GetByID retrieves a post by its ID
func (s *postService) GetByID(ctx context.Context, id uuid.UUID) (*models.PostResponse, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return s.mapPostToResponse(post), nil
}

// GetBySlug retrieves a post by its slug
func (s *postService) GetBySlug(ctx context.Context, slug string) (*models.PostResponse, error) {
	post, err := s.repo.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	return s.mapPostToResponse(post), nil
}

// GetAll retrieves all posts based on filter and pagination
func (s *postService) GetAll(ctx context.Context, filter *models.PostFilter, page, pageSize int) (*models.PaginatedPostResponse, error) {
	// Set default pagination values if not provided
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Calculate offset
	offset := (page - 1) * pageSize

	// Update filter with pagination
	if filter == nil {
		filter = &models.PostFilter{}
	}
	filter.Limit = pageSize
	filter.Offset = offset

	// Get posts and total count
	posts, total, err := s.repo.GetAll(filter)
	if err != nil {
		return nil, err
	}

	// Map posts to response
	postResponses := make([]*models.PostResponse, 0, len(posts))
	for _, post := range posts {
		postResponses = append(postResponses, s.mapPostToResponse(post))
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return &models.PaginatedPostResponse{
		Posts:      postResponses,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

// Update updates an existing post
func (s *postService) Update(ctx context.Context, id uuid.UUID, req *models.UpdatePostRequest) (*models.PostResponse, error) {
	// Get existing post
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// Check slug uniqueness if changed
	if req.Slug != "" && req.Slug != post.Slug {
		existingPost, err := s.repo.GetBySlug(req.Slug)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}
		if existingPost != nil && existingPost.ID != id {
			return nil, ErrPostSlugConflict
		}
	}

	// Update fields if provided
	if req.CategoryID != uuid.Nil {
		post.CategoryID = req.CategoryID
	}
	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Slug != "" {
		post.Slug = req.Slug
	}
	if req.Content != "" {
		post.Content = req.Content
	}
	if req.Excerpt != "" {
		post.Excerpt = req.Excerpt
	}
	if req.FeaturedImageURL != "" {
		post.FeaturedImageURL = req.FeaturedImageURL
	}
	if req.Status != "" {
		post.Status = req.Status
		// Update PublishedAt if status changes to published
		if req.Status == models.StatusPublished && post.PublishedAt == nil {
			now := time.Now()
			post.PublishedAt = &now
		}
	}
	if req.IsFeatured != nil {
		post.IsFeatured = *req.IsFeatured
	}
	if req.Metadata != nil {
		post.Metadata = req.Metadata
	}

	post.UpdatedAt = time.Now()

	// Update post and tags
	if err := s.repo.Update(post, req.TagIDs); err != nil {
		return nil, err
	}

	// Get updated post with all relationships
	updatedPost, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapPostToResponse(updatedPost), nil
}

// Delete soft-deletes a post
func (s *postService) Delete(ctx context.Context, id uuid.UUID) error {
	// Check if post exists
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrPostNotFound
		}
		return err
	}

	// Delete post
	return s.repo.Delete(post.ID)
}

// Publish changes a post's status to published and sets the published_at timestamp
func (s *postService) Publish(ctx context.Context, id uuid.UUID) (*models.PostResponse, error) {
	// Get existing post
	post, err := s.repo.GetByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrPostNotFound
		}
		return nil, err
	}

	// Set status to published and update published_at
	post.Status = models.StatusPublished
	now := time.Now()
	post.PublishedAt = &now
	post.UpdatedAt = now

	// Update post
	if err := s.repo.Update(post, nil); err != nil {
		return nil, err
	}

	// Get updated post
	updatedPost, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.mapPostToResponse(updatedPost), nil
}

// IncrementViewCount increments the view count of a post
func (s *postService) IncrementViewCount(ctx context.Context, id uuid.UUID) error {
	return s.repo.IncrementViewCount(id)
}

// mapPostToResponse maps a Post model to a PostResponse
func (s *postService) mapPostToResponse(post *models.Post) *models.PostResponse {
	if post == nil {
		return nil
	}

	// Parse metadata if exists
	var metadata interface{}
	if len(post.Metadata) > 0 {
		if err := json.Unmarshal(post.Metadata, &metadata); err != nil {
			// If unmarshal fails, use the raw bytes
			metadata = post.Metadata
		}
	}

	return &models.PostResponse{
		ID:               post.ID,
		AuthorID:         post.AuthorID,
		CategoryID:       post.CategoryID,
		Title:            post.Title,
		Slug:             post.Slug,
		Content:          post.Content,
		Excerpt:          post.Excerpt,
		FeaturedImageURL: post.FeaturedImageURL,
		Status:           post.Status,
		ViewCount:        post.ViewCount,
		IsFeatured:       post.IsFeatured,
		PublishedAt:      post.PublishedAt,
		CreatedAt:        post.CreatedAt,
		UpdatedAt:        post.UpdatedAt,
		Metadata:         metadata,
		Author:           post.Author,
		Category:         post.Category,
		Tags:             post.Tags,
	}
}
