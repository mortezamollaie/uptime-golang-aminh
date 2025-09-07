package models

import (
	"time"
)

type Node struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	URL       string    `gorm:"uniqueIndex;size:255" json:"url"`
	NodeLogs  []NodeLog `gorm:"foreignKey:NodeID" json:"node_logs"`
	Histories []History `gorm:"foreignKey:NodeID" json:"histories"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	// DeletedAt removed - check if table has this column
}

// TableName overrides the table name used by Node to `nodes`
func (Node) TableName() string {
	return "nodes"
}
