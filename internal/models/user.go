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
	Fullname     string    `json:"fullname" gorm:"type:varchar(255);default:''"`
	Email        string    `json:"email" gorm:"type:varchar(255);unique;not null"`
	Username     string    `json:"username" gorm:"type:varchar(255);unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"`
	Role         UserRole  `json:"role" gorm:"type:varchar(20);default:user"`
	AvatarURL    string    `json:"avatar_url"`
	IsActive     bool      `json:"is_active" gorm:"type:boolean;default:true"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty" gorm:"index"`

	Posts      []Post      `json:"posts,omitempty" gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	MediaFiles []MediaFile `json:"media_files,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	AuditLogs  []AuditLog  `json:"audit_logs,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RegisterRequest struct {
	Email    string   `json:"email" validate:"required,email"`
	Fullname string   `json:"fullname" validate:"required,min=2,max=100"`
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponse struct {
	User         UserResponse `json:"user"`
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	ExpiresIn    int64        `json:"expires_in"`
}

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Fullname  string    `json:"fullname"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Role      UserRole  `json:"role"`
	AvatarURL string    `json:"avatar_url"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
