package models

import (
	"gorm.io/gorm"
	"time"
)

type NodeLog struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	NodeID    uint           `gorm:"index" json:"node_id"` // کلید خارجی
	Delay     *float64       `json:"delay,omitempty"`
	Status    *uint          `json:"status,omitempty"`
	Up        bool           `gorm:"default:false" json:"up"`
	Suspended bool           `gorm:"default:false" json:"suspended"`
	Exception *string        `json:"exception,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
