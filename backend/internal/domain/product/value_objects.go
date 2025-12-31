package product

import "POSFlowBackend/internal/domain/shared"

type Category string

const (
	CategoryFood    Category = "Food"
	CategoryDrink   Category = "Drink"
	CategoryDessert Category = "Dessert"
)

func (c Category) IsValid() bool {
	switch c {
	case CategoryFood, CategoryDrink, CategoryDessert:
		return true
	default:
		return false
	}
}

type Stock struct {
	Quantity      int
	LowStockLevel int
}

func NewStock(quantity, lowStockLevel int) *Stock {
	return &Stock{
		Quantity:      quantity,
		LowStockLevel: lowStockLevel,
	}
}

func (s *Stock) IsLowStock() bool {
	return s.Quantity <= s.LowStockLevel
}

func (s *Stock) CanFulfill(quantity int) bool {
	return s.Quantity >= quantity
}

func (s *Stock) Decrease(quantity int) error {
	if !s.CanFulfill(quantity) {
		return shared.ErrInsufficientStock
	}
	s.Quantity -= quantity
	return nil
}

func (s *Stock) Increase(quantity int) {
	s.Quantity += quantity
}
