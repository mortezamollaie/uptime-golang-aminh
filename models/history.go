package models

import (
	"time"
)

type History struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NodeID    uint      `gorm:"index" json:"node_id"`
	Delay     *float64  `json:"delay,omitempty"`
	Status    *uint     `json:"status,omitempty"`
	Up        bool      `gorm:"default:false" json:"up"`
	Suspended bool      `gorm:"default:false" json:"suspended"`
	Exception *string   `json:"exception,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt removed - table doesn't have this column
}

// TableName overrides the table name used by History to `histories`
func (History) TableName() string {
	return "histories"
}
