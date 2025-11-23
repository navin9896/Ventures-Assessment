package models

import (
	"time"

	_ "github.com/jinzhu/gorm"
)

type User struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Password  string    `gorm:"not null" json:"-"`
	Token     *string   `gorm:"type:varchar(255)" json:"-"`
	CartID    *uint     `json:"cart_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relationships
	Cart   *Cart    `gorm:"foreignkey:CartID" json:"-"`
	Carts  []Cart   `gorm:"foreignkey:UserID" json:"-"`
	Orders []Order  `gorm:"foreignkey:UserID" json:"-"`
}

func (User) TableName() string {
	return "users"
}

