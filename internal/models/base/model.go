package base

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID        uuid.UUID       `json:"id" gorm:"type:uuid;default:uuid_generate_v7()"`
	CreatedAt *time.Time      `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt *time.Time      `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`
}
