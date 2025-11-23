package models

import (
	_ "github.com/jinzhu/gorm"
)

type CartItem struct {
	ID     uint `gorm:"primary_key" json:"id"`
	CartID uint `gorm:"not null" json:"cart_id"`
	ItemID uint `gorm:"not null" json:"item_id"`

	// Relationships
	Cart Cart `gorm:"foreignkey:CartID" json:"cart,omitempty"`
	Item Item `gorm:"foreignkey:ItemID" json:"item,omitempty"`
}

func (CartItem) TableName() string {
	return "cart_items"
}

