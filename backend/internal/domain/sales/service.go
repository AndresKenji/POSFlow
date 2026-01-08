package sales

import (
	"POSFlowBackend/internal/domain/order"
	"errors"

	"time"
)

// SalesService contains domain logic related to sales
type SalesService struct {
	salesRepo SalesRepository
	orderRepo order.OrderRepository
}

// DateRange represents a range between two dates
type DateRange struct {
	Start time.Time
	End   time.Time
}

func NewSalesService(salesRepo SalesRepository, orderRepo order.OrderRepository) *SalesService {
	return &SalesService{
		salesRepo: salesRepo,
		orderRepo: orderRepo,
	}
}

// CalculateDailySales calculates total sales for a given date
func (s *SalesService) CalculateDailySales(date time.Time) (*DailySales, error) {
	// Normalize date to midnight
	normalizedDate := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	// see if sales record already exists
	existing, err := s.salesRepo.FindByDate(normalizedDate)
	if err == nil {
		return existing, nil
	}

	// Create new daily sales record
	salesID := SalesID(normalizedDate.Format("2006-01-02"))
	dailySales := NewDailySales(salesID, normalizedDate)

	// Fetch orders for the day
	startOfDay := normalizedDate
	endOfDay := startOfDay.Add(24 * time.Hour)
	orders, err := s.orderRepo.FindByDateRange(startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}

	// add completed orders to daily sales
	for _, ord := range orders {
		if ord.IsCompleted() {
			dailySales.AddOrder(ord.ID(), ord.Total())
		}
	}

	// save daily sales record
	if err := s.salesRepo.Save(dailySales); err != nil {
		return nil, err
	}

	return dailySales, nil
}

// GetSalesReport generates a sales report for a given date range
func (s *SalesService) GetSalesReport(start, end time.Time) ([]*DailySales, error) {
	dateRange, err := NewDateRange(start, end)
	if err != nil {
		return nil, err
	}

	return s.salesRepo.FindByDateRange(dateRange.Start, dateRange.End)
}

func NewDateRange(start, end time.Time) (*DateRange, error) {
	if start.After(end) {
		return nil, errors.New("start date cannot be after end date")
	}
	return &DateRange{Start: start, End: end}, nil
}

// CloseDay closes the sales for a given day
func (s *SalesService) CloseDay(date time.Time) (*DailySales, error) {
	// calculate daily sales
	dailySales, err := s.CalculateDailySales(date)
	if err != nil {
		return nil, err
	}

	// close day
	if err := dailySales.CloseDay(); err != nil {
		return nil, err
	}

	// save updated daily sales
	if err := s.salesRepo.Save(dailySales); err != nil {
		return nil, err
	}

	return dailySales, nil
}
