package models

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	TotalAmount   float64           `json:"total_amount"`
	Status        string            `json:"status"`
	PaymentMethod string            `json:"payment_method"`
	CreatedAt     time.Time         `json:"created_at"`
	Items         []TransactionItem `json:"items,omitempty"`
}

type TransactionItem struct {
	ID            uuid.UUID `json:"id"`
	TransactionID uuid.UUID `json:"transaction_id"`
	BatteryID     uuid.UUID `json:"battery_id"`
	Quantity      int       `json:"quantity"`
	PriceAtTime   float64   `json:"price_at_time"`
	BatteryTitle  string    `json:"battery_title"`
	BatteryImg    string    `json:"battery_img"`
	CreatedAt     time.Time `json:"created_at"`
}
