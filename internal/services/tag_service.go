package services

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/repositories"
)

var (
	ErrTagNotFound     = errors.New("tag not found")
	ErrTagNameConflict = errors.New("tag name already exists")
	ErrTagSlugConflict = errors.New("tag slug already exists")
)

type TagService interface {
	Create(ctx context.Context, req *models.CreateTagRequest) (*models.TagResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.TagResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.TagResponse, error)
	GetAll(ctx context.Context, page, pageSize int) (*models.PaginatedTagResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateTagRequest) (*models.TagResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetTagsByPostID(ctx context.Context, postID uuid.UUID) ([]*models.TagResponse, error)
	AddTagsToPost(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error
	GetPostsByTagID(ctx context.Context, tagID uuid.UUID, page, pageSize int) (*models.PaginatedPostResponse, error)
}

type tagService struct {
	repo *repositories.TagRepository
}

func NewTagService(repo *repositories.TagRepository) TagService {
	return &tagService{
		repo: repo,
	}
}

func (s *tagService) Create(ctx context.Context, req *models.CreateTagRequest) (*models.TagResponse, error) {
	other, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if other != nil {
		return nil, ErrTagNameConflict
	}

	other, err = s.repo.GetBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if other != nil {
		return nil, ErrTagSlugConflict
	}

	tag := &models.Tag{
		Name:  req.Name,
		Slug:  req.Slug,
		Color: req.Color,
	}

	if err := s.repo.Create(tag); err != nil {
		return nil, err
	}

	return tag.ToResponse(), nil
}

func (s *tagService) GetByID(ctx context.Context, id uuid.UUID) (*models.TagResponse, error) {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag.ToResponse(), nil
}

func (s *tagService) GetBySlug(ctx context.Context, slug string) (*models.TagResponse, error) {
	tag, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag.ToResponse(), nil
}

func (s *tagService) GetAll(ctx context.Context, page, pageSize int) (*models.PaginatedTagResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	tags, total, err := s.repo.GetAll(pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	var responseTags []*models.TagResponse
	for _, t := range tags {
		responseTags = append(responseTags, t.ToResponse())
	}

	return &models.PaginatedTagResponse{
		Data:       responseTags,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *tagService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateTagRequest) (*models.TagResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTagNotFound
	}

	if req.Name != "" && req.Name != existing.Name {
		other, err := s.repo.GetByName(req.Name)
		if err != nil {
			return nil, err
		}
		if other != nil && other.ID != id {
			return nil, ErrTagNameConflict
		}
		existing.Name = req.Name
	}

	if req.Slug != "" && req.Slug != existing.Slug {
		other, err := s.repo.GetBySlug(req.Slug)
		if err != nil {
			return nil, err
		}
		if other != nil && other.ID != id {
			return nil, ErrTagSlugConflict
		}
		existing.Slug = req.Slug
	}

	if req.Color != "" {
		existing.Color = req.Color
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing.ToResponse(), nil
}

func (s *tagService) Delete(ctx context.Context, id uuid.UUID) error {
	tag, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if tag == nil {
		return ErrTagNotFound
	}

	return s.repo.Delete(id)
}

func (s *tagService) GetTagsByPostID(ctx context.Context, postID uuid.UUID) ([]*models.TagResponse, error) {
	tags, err := s.repo.GetTagsByPostID(postID)
	if err != nil {
		return nil, err
	}

	var responseTags []*models.TagResponse
	for _, t := range tags {
		responseTags = append(responseTags, t.ToResponse())
	}

	return responseTags, nil
}

func (s *tagService) AddTagsToPost(ctx context.Context, postID uuid.UUID, tagIDs []uuid.UUID) error {
	for _, tagID := range tagIDs {
		tag, err := s.repo.GetByID(tagID)
		if err != nil {
			return err
		}
		if tag == nil {
			return ErrTagNotFound
		}
	}

	return s.repo.AddTagsToPost(postID, tagIDs)
}

func (s *tagService) GetPostsByTagID(ctx context.Context, tagID uuid.UUID, page, pageSize int) (*models.PaginatedPostResponse, error) {
	tag, err := s.repo.GetByID(tagID)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	posts, total, err := s.repo.GetPostsByTagID(tagID, pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	var responsePosts []*models.PostResponse
	for _, p := range posts {
		// Map each post to a response using the same pattern as in post_service.go
		var metadata interface{}
		if len(p.Metadata) > 0 {
			if err := json.Unmarshal(p.Metadata, &metadata); err != nil {
				// If unmarshal fails, use the raw bytes
				metadata = p.Metadata
			}
		}
		
		responsePosts = append(responsePosts, &models.PostResponse{
			ID:               p.ID,
			AuthorID:         p.AuthorID,
			CategoryID:       p.CategoryID,
			Title:            p.Title,
			Slug:             p.Slug,
			Content:          p.Content,
			Excerpt:          p.Excerpt,
			FeaturedImageURL: p.FeaturedImageURL,
			Status:           p.Status,
			ViewCount:        p.ViewCount,
			IsFeatured:       p.IsFeatured,
			PublishedAt:      p.PublishedAt,
			CreatedAt:        p.CreatedAt,
			UpdatedAt:        p.UpdatedAt,
			Metadata:         metadata,
			Author:           p.Author,
			Category:         p.Category,
			Tags:             p.Tags,
		})
	}

	return &models.PaginatedPostResponse{
		Posts:      responsePosts,
		Total:      total, // Use int instead of int64
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}
