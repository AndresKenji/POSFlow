package commands

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type UpdateProductCommand struct {
	repo product.ProductRepository
}

func NewUpdateProductCommand(repo product.ProductRepository) *UpdateProductCommand {
	return &UpdateProductCommand{repo: repo}
}

func (c *UpdateProductCommand) Execute(id string, req dto.UpdateProductRequest) (*dto.ProductResponse, error) {
	// Find product
	prod, err := c.repo.FindByID(shared.ProductID(id))
	if err != nil {
		return nil, err
	}

	// Update price if provided
	if req.Price > 0 {
		price, err := shared.NewMoney(req.Price)
		if err != nil {
			return nil, err
		}
		if err := prod.UpdatePrice(*price); err != nil {
			return nil, err
		}
	}

	// Update info if provided
	if req.Name != "" || req.Description != "" || req.Category != "" {
		name := req.Name
		if name == "" {
			name = prod.Name()
		}

		description := req.Description
		if description == "" {
			description = prod.Description()
		}

		category := product.Category(req.Category)
		if req.Category == "" {
			category = prod.Category()
		}

		if err := prod.UpdateInfo(name, description, category); err != nil {
			return nil, err
		}
	}

	// Save changes
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
