package models

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

type User struct {
	ID           uuid.UUID `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	Email        string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Username     string    `json:"username" gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);default:user"`
	AvatarURL    string    `json:"avatar_url"`
	IsActive     bool      `json:"is_active" gorm:"type:boolean;default:true"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at" gorm:"index"`

	Posts      []Post      `json:"posts,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Comments   []Comment   `json:"comments,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MediaFiles []MediaFile `json:"media_files,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AuditLogs  []AuditLog  `json:"audit_logs,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
