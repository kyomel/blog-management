package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type AuditAction string

const (
	ActionCreate AuditAction = "create"
	ActionUpdate AuditAction = "update"
	ActionDelete AuditAction = "delete"
)

type AuditLog struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	UserID    uuid.UUID      `json:"user_id" gorm:"type:uuid;not null"`
	TableName string         `json:"table_name" gorm:"not null"`
	Action    AuditAction    `json:"action" gorm:"type:varchar(20);not null"`
	OldValues datatypes.JSON `json:"old_values" gorm:"type:jsonb"`
	NewValues datatypes.JSON `json:"new_values" gorm:"type:jsonb"`
	IPAddress string         `json:"ip_address"`
	UserAgent string         `json:"user_agent"`
	CreatedAt time.Time      `json:"created_at"`

	User User `json:"user" gorm:"foreignKey:UserID"`
}
