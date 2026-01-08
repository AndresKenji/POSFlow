package commands

import (
	"POSFlowBackend/internal/application/sales/dto"
	"POSFlowBackend/internal/domain/sales"
	"time"
)

type CloseDayCommand struct {
	salesService *sales.SalesService
}

func NewCloseDayCommand(salesService *sales.SalesService) *CloseDayCommand {
	return &CloseDayCommand{
		salesService: salesService,
	}
}

// Execute closes the sales for a given day
// If date is nil, it closes today's sales
func (c *CloseDayCommand) Execute(date *time.Time) (*dto.CloseDayResponse, error) {
	// if no date provided, use today's date
	targetDate := time.Now()
	if date != nil {
		targetDate = *date
	}

	// Close the day
	dailySales, err := c.salesService.CloseDay(targetDate)
	if err != nil {
		return &dto.CloseDayResponse{
			Success: false,
			Message: "Failed to close day: " + err.Error(),
		}, err
	}

	// Map to DTO
	salesResponse := &dto.DailySalesResponse{
		ID:          dailySales.ID().String(),
		Date:        dailySales.Date().Format("2006-01-02"),
		TotalSales:  dailySales.TotalSales().Amount,
		TotalOrders: dailySales.TotalOrders(),
		AverageSale: dailySales.AverageSale(),
		IsClosed:    dailySales.IsClosed(),
		ClosedAt:    dailySales.ClosedAt(),
		CreatedAt:   dailySales.CreatedAt(),
		UpdatedAt:   dailySales.UpdatedAt(),
	}

	return &dto.CloseDayResponse{
		Success:    true,
		Message:    "Day closed successfully",
		DailySales: salesResponse,
	}, nil
}
