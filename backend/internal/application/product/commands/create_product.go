package commands

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"

	"github.com/google/uuid"
)

type CreateProductCommand struct {
	repo product.ProductRepository
}

func NewCreateProductCommand(repo product.ProductRepository) *CreateProductCommand {
	return &CreateProductCommand{repo: repo}
}

func (c *CreateProductCommand) Execute(req dto.CreateProductRequest) (*dto.ProductResponse, error) {
	// Generate ID
	id := shared.ProductID(uuid.New().String())

	// Create Money value object
	price, err := shared.NewMoney(req.Price)
	if err != nil {
		return nil, err
	}

	// Map category
	category := product.Category(req.Category)

	// Create product entity using domain factory
	prod, err := product.NewProduct(id, req.Name, *price, category, req.Stock)
	if err != nil {
		return nil, err
	}

	// Add description if provided
	if req.Description != "" {
		prod.UpdateInfo(req.Name, req.Description, category)
	}

	// Save to repository
	if err := c.repo.Save(prod); err != nil {
		return nil, err
	}

	// Map to response DTO
	return &dto.ProductResponse{
		ID:          prod.ID().String(),
		Name:        prod.Name(),
		Description: prod.Description(),
		Price:       prod.Price().Amount,
		Category:    string(prod.Category()),
		Stock:       prod.Stock(),
		Active:      prod.IsActive(),
		IsLowStock:  prod.IsLowStock(),
		CreatedAt:   prod.CreatedAt(),
		UpdatedAt:   prod.UpdatedAt(),
	}, nil
}
