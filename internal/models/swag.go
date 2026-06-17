package models

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity" validate:"required,min=1"`
}

type AddToFavoritesRequest struct {
	BatteryID string `json:"battery_id" validate:"required"`
}
