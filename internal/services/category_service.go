package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/kyomel/blog-management/internal/models"
	"github.com/kyomel/blog-management/internal/repositories"
)

var (
	ErrCategoryNotFound     = errors.New("category not found")
	ErrCategoryNameConflict = errors.New("category name already exists")
	ErrCategorySlugConflict = errors.New("category slug already exists")
)

type CategoryService interface {
	Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.CategoryResponse, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.CategoryResponse, error)
	GetBySlug(ctx context.Context, slug string) (*models.CategoryResponse, error)
	GetAll(ctx context.Context, page, pageSize int) (*models.PaginatedCategoryResponse, error)
	Update(ctx context.Context, id uuid.UUID, req *models.UpdateCategoryRequest) (*models.CategoryResponse, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type categoryService struct {
	repo *repositories.CategoryRepository
}

func NewCategoryService(repo *repositories.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) Create(ctx context.Context, req *models.CreateCategoryRequest) (*models.CategoryResponse, error) {
	other, err := s.repo.GetByName(req.Name)
	if err != nil {
		return nil, err
	}
	if other != nil {
		return nil, ErrCategoryNameConflict
	}

	other, err = s.repo.GetBySlug(req.Slug)
	if err != nil {
		return nil, err
	}
	if other != nil {
		return nil, ErrCategorySlugConflict
	}

	category := &models.Category{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
	}

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}

	return category.ToResponse(), nil
}

func (s *categoryService) GetByID(ctx context.Context, id uuid.UUID) (*models.CategoryResponse, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category.ToResponse(), nil
}

func (s *categoryService) GetBySlug(ctx context.Context, slug string) (*models.CategoryResponse, error) {
	category, err := s.repo.GetBySlug(slug)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}
	return category.ToResponse(), nil
}

func (s *categoryService) GetAll(ctx context.Context, page, pageSize int) (*models.PaginatedCategoryResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	categories, total, err := s.repo.GetAll(pageSize, offset)
	if err != nil {
		return nil, err
	}

	totalPages := (total + pageSize - 1) / pageSize
	if totalPages == 0 {
		totalPages = 1
	}

	var responseCategories []*models.CategoryResponse
	for _, c := range categories {
		responseCategories = append(responseCategories, c.ToResponse())
	}

	return &models.PaginatedCategoryResponse{
		Data:       responseCategories,
		Total:      int64(total),
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}, nil
}

func (s *categoryService) Update(ctx context.Context, id uuid.UUID, req *models.UpdateCategoryRequest) (*models.CategoryResponse, error) {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrCategoryNotFound
	}

	if req.Name != "" && req.Name != existing.Name {
		other, err := s.repo.GetByName(req.Name)
		if err != nil {
			return nil, err
		}
		if other != nil && other.ID != id {
			return nil, ErrCategoryNameConflict
		}
		existing.Name = req.Name
	}

	if req.Slug != "" && req.Slug != existing.Slug {
		other, err := s.repo.GetBySlug(req.Slug)
		if err != nil {
			return nil, err // DB error
		}
		if other != nil && other.ID != id {
			return nil, ErrCategorySlugConflict
		}
		existing.Slug = req.Slug
	}

	if req.Description != "" {
		existing.Description = req.Description
	}

	if err := s.repo.Update(existing); err != nil {
		return nil, err
	}

	return existing.ToResponse(), nil
}

func (s *categoryService) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}
