package models

import "gorm.io/gorm"

type Item struct {
	gorm.Model
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	SellerID    uint    `json:"seller_id"`
	Seller      User    `json:"seller,omitempty"`
}
