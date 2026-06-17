package models

import (
	"time"

	"github.com/google/uuid"
)

type CartItem struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	BatteryID   uuid.UUID `json:"battery_id"`
	Quantity    int       `json:"quantity"`
	BatteryName string    `json:"battery_name"`
	Price       float64   `json:"price"`
	Img         string    `json:"img"`
	Brand       string    `json:"brand"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AddToCartRequest struct {
	BatteryID uuid.UUID `json:"battery_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}
