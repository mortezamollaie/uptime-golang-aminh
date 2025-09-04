package models

import (
	"time"

	"gorm.io/gorm"
)

type Node struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	URL       string         `gorm:"uniqueIndex;size:255" json:"url"`
	NodeLogs  []NodeLog      `gorm:"foreignKey:NodeID" json:"node_logs"`
	Histories []History      `gorm:"foreignKey:NodeID" json:"histories"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
