package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ToMonitor struct {
	ID        uuid.UUID `gorm:"type:varchar(36);primaryKey;"`
	Site      string    `gorm:"unique" json:"site"`
	Webhook   string    `json:"webhook"`
	Interval  int       `json:"interval"`
	CreatedAt time.Time
}

func (r *ToMonitor) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}
