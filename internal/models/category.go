package models

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Slug        string    `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Posts []Post `json:"posts,omitempty" gorm:"foreignKey:CategoryID"`
}

type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=255"`
	Slug        string `json:"slug" validate:"required,slug,max=255"`
	Description string `json:"description,omitempty"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
	Slug        string `json:"slug,omitempty" validate:"omitempty,slug,max=255"`
	Description string `json:"description,omitempty"`
}

type CategoryResponse struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (c *Category) ToResponse() *CategoryResponse {
	return &CategoryResponse{
		ID:          c.ID,
		Name:        c.Name,
		Slug:        c.Slug,
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

type PaginatedCategoryResponse struct {
	Data       []*CategoryResponse `json:"data"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PageSize   int                 `json:"page_size"`
	TotalPages int                 `json:"total_pages"`
}
