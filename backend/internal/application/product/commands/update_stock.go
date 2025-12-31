package commands

import (
	"POSFlowBackend/internal/application/product/dto"
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type UpdateStockCommand struct {
	repo product.ProductRepository
}

func NewUpdateStockCommand(repo product.ProductRepository) *UpdateStockCommand {
	return &UpdateStockCommand{repo: repo}
}

func (c *UpdateStockCommand) Execute(id string, req dto.UpdateStockRequest) (*dto.ProductResponse, error) {
	// Find product
	prod, err := c.repo.FindByID(shared.ProductID(id))
	if err != nil {
		return nil, err
	}

	// Update stock based on type
	switch req.Type {
	case "add":
		if err := prod.IncreaseStock(req.Quantity); err != nil {
			return nil, err
		}
	case "remove":
		if err := prod.DecreaseStock(req.Quantity); err != nil {
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
