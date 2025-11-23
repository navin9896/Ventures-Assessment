package models

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

type Cart struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	UserID    uint      `gorm:"not null" json:"user_id"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	User      User       `gorm:"foreignkey:UserID" json:"user,omitempty"`
	CartItems []CartItem `gorm:"foreignkey:CartID" json:"cart_items,omitempty"`
	Orders    []Order    `gorm:"foreignkey:CartID" json:"-"`
}

func (Cart) TableName() string {
	return "carts"
}

