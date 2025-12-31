package dto

import "time"

// CreateOrderRequest - Input DTO for creating an order
type CreateOrderRequest struct {
	TableNumber string      `json:"table_number" binding:"required"`
	Items       []OrderItem `json:"items" binding:"required,min=1"`
}

type OrderItem struct {
	ProductID string `json:"product_id" binding:"required"`
	Quantity  int    `json:"quantity" binding:"required,gt=0"`
}

// UpdateOrderStatusRequest - Input DTO for updating order status
type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=pending preparing ready completed cancelled"`
}

// OrderResponse - Output DTO
type OrderResponse struct {
	ID          string              `json:"id"`
	TableNumber string              `json:"table_number"`
	Status      string              `json:"status"`
	Items       []OrderItemResponse `json:"items"`
	Total       float64             `json:"total"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type OrderItemResponse struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	UnitPrice   float64 `json:"unit_price"`
	Subtotal    float64 `json:"subtotal"`
}

// OrderListResponse - Output DTO for list
type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  int              `json:"total"`
}
