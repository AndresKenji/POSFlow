package order

import (
	"POSFlowBackend/internal/domain/shared"
	"time"
)

// OrderRepository defines the interface for order persistence
type OrderRepository interface {
	Save(order *Order) error
	FindByID(id shared.OrderID) (*Order, error)
	FindAll() ([]*Order, error)
	FindPending() ([]*Order, error)
	FindByStatus(status OrderStatus) ([]*Order, error)
	FindByDateRange(start, end time.Time) ([]*Order, error)
}
