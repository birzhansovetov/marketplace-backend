package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ItemID  uint    `json:"item_id"`
	Item    Item    `json:"item,omitempty"`
	BuyerID uint    `json:"buyer_id"`
	Buyer   User    `json:"buyer,omitempty"`
	Status  string  `json:"status"` // pending, confirmed, cancelled, completed
	Price   float64 `json:"price"`
}
