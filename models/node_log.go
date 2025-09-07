package models

import (
	"time"
)

type NodeLog struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	NodeID    uint      `gorm:"index" json:"node_id"` // کلید خارجی
	Delay     *float64  `json:"delay,omitempty"`
	Status    *uint     `json:"status,omitempty"`
	Up        bool      `gorm:"default:false" json:"up"`
	Suspended bool      `gorm:"default:false" json:"suspended"`
	Exception *string   `json:"exception,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt removed - check if table has this column
}

// TableName overrides the table name used by NodeLog to `node_logs`
func (NodeLog) TableName() string {
	return "node_logs"
}
