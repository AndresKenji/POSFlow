package product

import "POSFlowBackend/internal/domain/shared"

// ProductRepository defines the interface for product persistence
// Implementation will be in infrastructure layer
type ProductRepository interface {
	Save(product *Product) error
	FindByID(id shared.ProductID) (*Product, error)
	FindAll() ([]*Product, error)
	FindByCategory(category Category) ([]*Product, error)
	FindLowStock() ([]*Product, error)
	Delete(id shared.ProductID) error
}
