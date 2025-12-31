package order

import (
	"POSFlowBackend/internal/domain/product"
	"POSFlowBackend/internal/domain/shared"
)

// OrderService contains domain logic for orders
type OrderService struct {
	orderRepo   OrderRepository
	productRepo product.ProductRepository
}

func NewOrderService(orderRepo OrderRepository, productRepo product.ProductRepository) *OrderService {
	return &OrderService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
	}
}

// CreateOrder handles order creation with stock validation
func (s *OrderService) CreateOrder(
	id shared.OrderID,
	tableNumber TableNumber,
	itemRequests []struct {
		ProductID shared.ProductID
		Quantity  int
	},
) (*Order, error) {

	var orderItems []*OrderItem

	// Validate stock and create order items
	for _, req := range itemRequests {
		product, err := s.productRepo.FindByID(req.ProductID)
		if err != nil {
			return nil, err
		}

		if !product.IsActive() {
			return nil, shared.ErrInvalidInput
		}

		// Check stock availability
		if err := product.DecreaseStock(req.Quantity); err != nil {
			return nil, err
		}

		// Create order item
		item, err := NewOrderItem(req.ProductID, req.Quantity, product.Price())
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, item)

		// Save updated product stock
		if err := s.productRepo.Save(product); err != nil {
			return nil, err
		}
	}

	// Create order
	order, err := NewOrder(id, tableNumber, orderItems)
	if err != nil {
		return nil, err
	}

	// Save order
	if err := s.orderRepo.Save(order); err != nil {
		return nil, err
	}

	return order, nil
}

// GetPendingOrders returns orders that need attention in the kitchen
func (s *OrderService) GetPendingOrders() ([]*Order, error) {
	return s.orderRepo.FindPending()
}
