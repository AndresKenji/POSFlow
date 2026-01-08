package sales

import (
	"POSFlowBackend/internal/domain/shared"
	"time"
)

type DailySales struct {
	id          SalesID
	date        time.Time
	totalSales  shared.Money
	totalOrders int
	orderIDs    []shared.OrderID
	closed      bool
	closedAt    *time.Time
	createdAt   time.Time
	updatedAt   time.Time
}

func NewDailySales(id SalesID, date time.Time) *DailySales {
	return &DailySales{
		id:          id,
		date:        date,
		totalSales:  shared.Money{Amount: 0, Currency: "USD"},
		totalOrders: 0,
		orderIDs:    []shared.OrderID{},
		closed:      false,
		closedAt:    nil,
		createdAt:   time.Now(),
		updatedAt:   time.Now(),
	}
}

// Getters
func (s *DailySales) ID() SalesID                { return s.id }
func (s *DailySales) Date() time.Time            { return s.date }
func (s *DailySales) TotalSales() shared.Money   { return s.totalSales }
func (s *DailySales) TotalOrders() int           { return s.totalOrders }
func (s *DailySales) OrderIDs() []shared.OrderID { return s.orderIDs }
func (s *DailySales) IsClosed() bool             { return s.closed }
func (s *DailySales) ClosedAt() *time.Time       { return s.closedAt }
func (s *DailySales) CreatedAt() time.Time       { return s.createdAt }
func (s *DailySales) UpdatedAt() time.Time       { return s.updatedAt }

// Business Logic
func (s *DailySales) AddOrder(orderID shared.OrderID, amount shared.Money) error {
	if s.closed {
		return shared.ErrInvalidInput // Day is closed for adding orders
	}

	s.orderIDs = append(s.orderIDs, orderID)
	s.totalSales = s.totalSales.Add(amount)
	s.totalOrders++
	s.updatedAt = time.Now()

	return nil
}

func (s *DailySales) CloseDay() error {
	if s.closed {
		return shared.ErrInvalidInput // Day is already closed
	}

	now := time.Now()
	s.closed = true
	s.closedAt = &now
	s.updatedAt = time.Now()

	return nil
}

func (s *DailySales) AverageSale() float64 {
	if s.totalOrders == 0 {
		return 0
	}
	return s.totalSales.Amount / float64(s.totalOrders)
}

func ReconstructDailySales(
	id SalesID,
	date time.Time,
	totalSales shared.Money,
	totalOrders int,
	orderIDs []shared.OrderID,
	closed bool,
	closedAt *time.Time,
	createdAt time.Time,
	updatedAt time.Time,
) *DailySales {
	return &DailySales{
		id:          id,
		date:        date,
		totalSales:  totalSales,
		totalOrders: totalOrders,
		orderIDs:    orderIDs,
		closed:      closed,
		closedAt:    closedAt,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}
}
