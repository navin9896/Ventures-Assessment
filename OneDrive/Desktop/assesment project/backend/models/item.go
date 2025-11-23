package models

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

type Item struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	CartItems []CartItem `gorm:"foreignkey:ItemID" json:"-"`
}

func (Item) TableName() string {
	return "items"
}

