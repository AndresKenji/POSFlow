package sales

import "time"

// SalesRepository defines the methods for persisting and retrieving sales data
type SalesRepository interface {
	// Save save daily sales record to the database
	Save(sales *DailySales) error

	// Find daily sales by its ID
	FindByID(id SalesID) (*DailySales, error)

	// Find daily sales by date
	FindByDate(date time.Time) (*DailySales, error)

	// Find daily sales within a date range
	FindByDateRange(start, end time.Time) ([]*DailySales, error)

	// Get all daily sales records
	FindAll() ([]*DailySales, error)

	// Verify if a daily sales record exists for a given date
	ExistsByDate(date time.Time) (bool, error)
}
