package models

import "time"

type DomainCheckHistory struct {
	ID        int64     `gorm:"column:id" json:"id"`
	Name      string    `gorm:"column:name" json:"name"`
	Total     int       `gorm:"column:total" json:"total"`
	Success   int       `gorm:"column:success" json:"success"`
	Failed    int       `gorm:"column:failed" json:"failed"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`

	Domains []Domain `gorm:"foreignKey:BatchId;references:ID"`
}
