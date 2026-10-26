package model

import "time"

type OrderStatus string

const (
	StatusCreated   OrderStatus = "CREATED"
	StatusConfirmed OrderStatus = "CONFIRMED"
	StatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID         int64       `json:"id"`
	CustomerID int         `json:"customerId"`
	ProductID  int64       `json:"productId"`
	Quantity   int         `json:"quantity"`
	Total      float64     `json:"total"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"createdAt"`
}

type CreateOrderRequest struct {
	CustomerID int   `json:"customerId" binding:"required,gt=0"`
	ProductID  int64 `json:"productId" binding:"required,gt=0"`
	Quantity   int   `json:"quantity" binding:"required,gt=0"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required,oneof=CREATED CONFIRMED CANCELLED"`
}
