package dto

import (
	"time"
)

// CreateProductRequest - Input DTO for creating a product
type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Category    string  `json:"category" binding:"required"`
	Stock       int     `json:"stock" binding:"gte=0"`
}

// UpdateProductRequest - Input DTO for updating a product
type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"gt=0"`
	Category    string  `json:"category"`
}

// UpdateStockRequest - Input DTO for updating stock
type UpdateStockRequest struct {
	Quantity int    `json:"quantity" binding:"required,gt=0"`
	Type     string `json:"type" binding:"required,oneof=add remove"`
}

// ProductResponse - Output DTO
type ProductResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Category    string    `json:"category"`
	Stock       int       `json:"stock"`
	Active      bool      `json:"active"`
	IsLowStock  bool      `json:"is_low_stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ProductListResponse - Output DTO for list
type ProductListResponse struct {
	Products []*ProductResponse `json:"products"`
	Total    int                `json:"total"`
}
