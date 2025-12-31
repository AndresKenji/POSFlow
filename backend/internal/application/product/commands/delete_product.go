package commands

import (
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

type DeleteProductCommand struct {
	repo product.ProductRepository
}

func NewDeleteProductCommand(repo product.ProductRepository) *DeleteProductCommand {
	return &DeleteProductCommand{repo: repo}
}

func (c *DeleteProductCommand) Execute(id string) error {
	// Find product first to ensure it exists
	prod, err := c.repo.FindByID(shared.ProductID(id))
	if err != nil {
		return err
	}

	// Soft delete by deactivating
	prod.Deactivate()

	// Save changes
	if err := c.repo.Save(prod); err != nil {
		return err
	}

	return nil
}
