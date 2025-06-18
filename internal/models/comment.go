package models

import (
	"time"

	"github.com/google/uuid"
)

type CommentStatus string

const (
	CommentStatusPending  CommentStatus = "pending"
	CommentStatusApproved CommentStatus = "approved"
	CommentStatusRejected CommentStatus = "rejected"
)

type Comment struct {
	ID        uuid.UUID     `json:"id" gorm:"type:uuid;primarykey;default:gen_random_uuid()"`
	PostID    uuid.UUID     `json:"post_id" gorm:"type:uuid;not null"`
	UserID    uuid.UUID     `json:"user_id" gorm:"type:uuid;not null"`
	ParentID  *uuid.UUID    `json:"parent_id" gorm:"type:uuid"`
	Content   string        `json:"content" gorm:"type:text;not null"`
	Status    CommentStatus `json:"status" gorm:"type:varchar(20);default:pending"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	DeletedAt time.Time     `json:"deleted_at" gorm:"index"`

	Post    *Post     `json:"post,omitempty" gorm:"foreignKey:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User    *User     `json:"user,omitempty" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Parent  *Comment  `json:"parent,omitempty" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Replies []Comment `json:"replies,omitempty" gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
