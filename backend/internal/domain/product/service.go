package product

import "POSFlowBackend/internal/domain/shared"

// ProductService contains domain logic that doesn't fit in entities
type ProductService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) *ProductService {
	return &ProductService{repo: repo}
}

// CheckStockAvailability checks if we have enough stock for an order
func (s *ProductService) CheckStockAvailability(productID shared.ProductID, quantity int) error {
	product, err := s.repo.FindByID(productID)
	if err != nil {
		return err
	}

	if !product.IsActive() {
		return shared.ErrInvalidInput
	}

	if !product.stock.CanFulfill(quantity) {
		return shared.ErrInsufficientStock
	}

	return nil
}

// GetLowStockProducts returns products that need restocking
func (s *ProductService) GetLowStockProducts() ([]*Product, error) {
	return s.repo.FindLowStock()
}
