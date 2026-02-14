package models

import "time"

type Domain struct {
	Id          int64     `gorm:"column:id" json:"id"`
	BatchId     int64     `gorm:"column:batch_id" json:"batch_id"`
	FullUrl     string    `gorm:"column:full_url" json:"full_url"`
	StatusCode  int       `gorm:"column:status_code" json:"status_code"`
	RedirectUrl string    `gorm:"column:redirect_url" json:"redirect_url"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at" json:"updated_at"`
}
