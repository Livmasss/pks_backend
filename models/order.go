package models

import "time"

type Order struct {
	OrderID   int        `db:"order_id" json:"order_id"`
	UserID    string     `db:"user_id" json:"user_id"`
	Total     float64    `db:"total" json:"total"`
	Status    string     `db:"status" json:"status"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	CartItems []CartItem `json:"products"`
}
