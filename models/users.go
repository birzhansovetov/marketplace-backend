package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
	Items    []Item `json:"items,omitempty" gorm:"foreignKey:SellerID"`
}
