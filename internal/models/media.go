package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MediaFile struct {
	ID                 uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID             uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	OriginalName       string         `json:"original_name" gorm:"not null"`
	FileName           string         `json:"file_name" gorm:"not null"`
	FilePath           string         `json:"file_path" gorm:"not null"`
	CloudinaryPublicID string         `json:"cloudinary_public_id"`
	MimeType           string         `json:"mime_type" gorm:"not null"`
	FileSize           int64          `json:"file_size" gorm:"not null"`
	Metadata           datatypes.JSON `json:"metadata" gorm:"type:jsonb"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at" gorm:"index"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}
