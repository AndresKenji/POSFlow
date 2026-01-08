package queries

import (
	"POSFlowBackend/internal/application/sales/dto"
	"POSFlowBackend/internal/domain/sales"
	"time"
)

type GetSalesReportQuery struct {
	salesService *sales.SalesService
}

func NewGetSalesReportQuery(salesService *sales.SalesService) *GetSalesReportQuery {
	return &GetSalesReportQuery{
		salesService: salesService,
	}
}

// Execute obtains the sales report for a given date range
func (q *GetSalesReportQuery) Execute(startDate, endDate time.Time) (*dto.SalesReportResponse, error) {
	// get daily sales within the date range
	dailySalesList, err := q.salesService.GetSalesReport(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Map to DTO and calculate totals
	var totalSales float64
	var totalOrders int
	var dailySalesResponses []*dto.DailySalesResponse

	for _, ds := range dailySalesList {
		totalSales += ds.TotalSales().Amount
		totalOrders += ds.TotalOrders()

		dailySalesResponses = append(dailySalesResponses, &dto.DailySalesResponse{
			ID:          ds.ID().String(),
			Date:        ds.Date().Format("2006-01-02"),
			TotalSales:  ds.TotalSales().Amount,
			TotalOrders: ds.TotalOrders(),
			AverageSale: ds.AverageSale(),
			IsClosed:    ds.IsClosed(),
			ClosedAt:    ds.ClosedAt(),
			CreatedAt:   ds.CreatedAt(),
			UpdatedAt:   ds.UpdatedAt(),
		})
	}

	// Calculate average sale
	averageSale := 0.0
	if totalOrders > 0 {
		averageSale = totalSales / float64(totalOrders)
	}

	return &dto.SalesReportResponse{
		StartDate:   startDate.Format("2006-01-02"),
		EndDate:     endDate.Format("2006-01-02"),
		TotalSales:  totalSales,
		TotalOrders: totalOrders,
		DailySales:  dailySalesResponses,
		AverageSale: averageSale,
	}, nil
}
