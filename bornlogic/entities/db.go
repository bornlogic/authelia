package entities

import "time"

type DB struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	CreatedTime time.Time `json:"created_at"`
	UpdatedTime time.Time `json:"updated_at"`
}
