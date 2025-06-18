package models

import (
	"time"

	"github.com/google/uuid"
)

type PostFilter struct {
	Status     PostStatus
	CategoryID *uuid.UUID
	AuthorID   *uuid.UUID
	IsFeatured *bool
	Search     string
	Limit      int
	Offset     int
}

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
	StatusArchived  PostStatus = "archived"
)

type CreatePostRequest struct {
	AuthorID         uuid.UUID   `json:"author_id" validate:"required"`
	CategoryID       uuid.UUID   `json:"category_id" validate:"required"`
	Title            string      `json:"title" validate:"required"`
	Slug             string      `json:"slug" validate:"required"`
	Content          string      `json:"content" validate:"required"`
	Excerpt          string      `json:"excerpt"`
	FeaturedImageURL string      `json:"featured_image_url"`
	Status           PostStatus  `json:"status" validate:"required,oneof=draft published archived"`
	IsFeatured       bool        `json:"is_featured"`
	Metadata         []byte      `json:"metadata,omitempty"`
	TagIDs           []uuid.UUID `json:"tag_ids,omitempty"`
}

type UpdatePostRequest struct {
	CategoryID       uuid.UUID   `json:"category_id,omitempty"`
	Title            string      `json:"title,omitempty"`
	Slug             string      `json:"slug,omitempty"`
	Content          string      `json:"content,omitempty"`
	Excerpt          string      `json:"excerpt,omitempty"`
	FeaturedImageURL string      `json:"featured_image_url,omitempty"`
	Status           PostStatus  `json:"status,omitempty" validate:"omitempty,oneof=draft published archived"`
	IsFeatured       *bool       `json:"is_featured,omitempty"`
	Metadata         []byte      `json:"metadata,omitempty"`
	TagIDs           []uuid.UUID `json:"tag_ids,omitempty"`
}

// PostResponse represents the response for a post
type PostResponse struct {
	ID               uuid.UUID   `json:"id"`
	AuthorID         uuid.UUID   `json:"author_id"`
	CategoryID       uuid.UUID   `json:"category_id"`
	Title            string      `json:"title"`
	Slug             string      `json:"slug"`
	Content          string      `json:"content"`
	Excerpt          string      `json:"excerpt"`
	FeaturedImageURL string      `json:"featured_image_url"`
	Status           PostStatus  `json:"status"`
	ViewCount        int         `json:"view_count"`
	IsFeatured       bool        `json:"is_featured"`
	PublishedAt      *time.Time  `json:"published_at,omitempty"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
	Metadata         interface{} `json:"metadata,omitempty"`

	Author   *User     `json:"author,omitempty"`
	Category *Category `json:"category,omitempty"`
	Tags     []*Tag    `json:"tags,omitempty"`
}

type PaginatedPostResponse struct {
	Posts      []*PostResponse `json:"posts"`
	Total      int             `json:"total"`
	Page       int             `json:"page"`
	PageSize   int             `json:"page_size"`
	TotalPages int             `json:"total_pages"`
}

type Post struct {
	ID               uuid.UUID  `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	AuthorID         uuid.UUID  `json:"author_id" gorm:"type:uuid;not null"`
	CategoryID       uuid.UUID  `json:"category_id" gorm:"type:uuid;not null"`
	Title            string     `json:"title" gorm:"type:varchar(255);not null"`
	Slug             string     `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Content          string     `json:"content" gorm:"type:text"`
	Excerpt          string     `json:"excerpt" gorm:"type:text"`
	FeaturedImageURL string     `json:"featured_image_url"`
	Status           PostStatus `json:"status" gorm:"type:varchar(20);default:draft"`
	ViewCount        int        `json:"view_count" gorm:"type:int;default:0"`
	IsFeatured       bool       `json:"is_featured" gorm:"type:boolean;default:false"`
	PublishedAt      *time.Time `json:"published_at"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
	DeletedAt        *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Metadata         []byte     `json:"metadata,omitempty"`

	Author   *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tags     []*Tag    `json:"tags,omitempty" gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
