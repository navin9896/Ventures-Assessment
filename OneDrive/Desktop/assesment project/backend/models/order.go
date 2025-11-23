package models

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

type Order struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CartID    uint      `gorm:"not null" json:"cart_id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Cart Cart `gorm:"foreignkey:CartID" json:"cart,omitempty"`
	User User `gorm:"foreignkey:UserID" json:"user,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}

