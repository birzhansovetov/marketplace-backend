package dto

type CreateOrderRequest struct {
	ItemID  uint `json:"item_id" binding:"required"`
	BuyerID uint `json:"buyer_id" binding:"required"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending confirmed cancelled completed"`
}
