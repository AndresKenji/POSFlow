package queries

import (
	"POSFlowBackend/internal/application/sales/dto"
	"POSFlowBackend/internal/domain/sales"
	"time"
)

type GetDailySalesQuery struct {
	salesService *sales.SalesService
}

func NewGetDailySalesQuery(salesService *sales.SalesService) *GetDailySalesQuery {
	return &GetDailySalesQuery{
		salesService: salesService,
	}
}

func (q *GetDailySalesQuery) Execute(date *time.Time) (*dto.DailySalesResponse, error) {
	targetDate := time.Now()
	if date != nil {
		targetDate = *date
	}

	dailySales, err := q.salesService.CalculateDailySales(targetDate)
	if err != nil {
		return nil, err
	}

	return q.mapToDTO(dailySales), nil
}

func (q *GetDailySalesQuery) mapToDTO(ds *sales.DailySales) *dto.DailySalesResponse {
	return &dto.DailySalesResponse{
		ID:          ds.ID().String(),
		Date:        ds.Date().Format("2006-01-02"),
		TotalSales:  ds.TotalSales().Amount,
		TotalOrders: ds.TotalOrders(),
		AverageSale: ds.AverageSale(),
		IsClosed:    ds.IsClosed(),
		ClosedAt:    ds.ClosedAt(),
		CreatedAt:   ds.CreatedAt(),
		UpdatedAt:   ds.UpdatedAt(),
	}
}
