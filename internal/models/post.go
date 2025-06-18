package models

import (
	"time"

	"github.com/google/uuid"
)

type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
	StatusArchived  PostStatus = "archived"
)

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
	DeletedAt        time.Time  `json:"deleted_at" gorm:"index"`

	Author   *User     `json:"author,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Category *Category `json:"category,omitempty" gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments []Comment `json:"comments,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Tags     []Tag     `json:"tags,omitempty" gorm:"many2many:post_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
