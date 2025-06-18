package models

import (
	"time"

	"github.com/google/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Name      string    `json:"name" gorm:"type:varchar(255);uniqueIndex;not null"`
	Slug      string    `json:"slug" gorm:"type:varchar(255);uniqueIndex;not null"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	DeletedAt time.Time `json:"deleted_at" gorm:"index"`

	Posts []Post `json:"posts,omitempty" gorm:"many2many:post_tags;"`
}
