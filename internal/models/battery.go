package models

import (
	"github.com/google/uuid"
	"time"
)

type Battery struct {
	ID         uuid.UUID `json:"id"`
	Title      string    `json:"title"`
	Price      float64   `json:"price"`
	Stock      int       `json:"stock"`
	Img        string    `json:"img"`
	Brand      string    `json:"brand"`
	Voltage    int       `json:"voltage"`
	Polarity   string    `json:"polarity"`
	Capacity   float64   `json:"capacity"`
	Standart   string    `json:"standart"`
	Technology string    `json:"technology"`
	SizeType   string    `json:"sizeType"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
